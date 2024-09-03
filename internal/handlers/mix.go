package handlers

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gabe565/relax-sounds/internal/encoder/encode"
	"github.com/gabe565/relax-sounds/internal/encoder/filetype"
	"github.com/gabe565/relax-sounds/internal/preset"
	"github.com/gabe565/relax-sounds/internal/stream"
	"github.com/gabe565/relax-sounds/internal/stream/streamcache"
	"github.com/gopxl/beep"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/spf13/cobra"
)

func MixFlags(cmd *cobra.Command) {
	chunkLengthDefault := 2 * time.Minute
	if env := os.Getenv("STREAM_CHUNK_LENGTH"); env != "" {
		var err error
		chunkLengthDefault, err = time.ParseDuration(env)
		if err != nil {
			slog.Warn("Failed to parse STREAM_CHUNK_LENGTH env", "error", err)
		}
	}
	cmd.PersistentFlags().Duration("stream-chunk-length", chunkLengthDefault, "Sets the length of each chunk when casting")
}

//nolint:gocyclo
func Mix(app *pocketbase.PocketBase) echo.HandlerFunc {
	cache := streamcache.New()
	dataFs := os.DirFS(filepath.Join(app.DataDir(), "storage"))
	chunkLength, err := app.RootCmd.PersistentFlags().GetDuration("stream-chunk-length")
	if err != nil {
		panic(err)
	}

	return func(c echo.Context) error {
		var err error

		query := c.PathParam("query")

		// Preset parameter
		presetEncoded, fileTypeStr, found := strings.Cut(query, ".")
		if !found {
			return apis.NewNotFoundError("", nil)
		}

		// File type parameter
		var fileType filetype.FileType
		err = fileType.UnmarshalText([]byte(fileTypeStr))
		if err != nil {
			if errors.Is(err, filetype.ErrInvalidFileType) {
				return apis.NewNotFoundError("", nil)
			}
			panic(err)
		}

		uuid := c.PathParam("uuid")
		entry, found := cache.Get(uuid)
		if found && entry.Preset != presetEncoded {
			// Same window is changing streams
			// Destroy old stream then recreate
			entry.Log.Info("Close stream",
				"accessed", entry.Accessed,
				"age", time.Since(entry.Created).Round(100*time.Millisecond).String(),
				"transferred", humanize.IBytes(entry.Transferred),
			)
			if err := entry.Close(); err != nil {
				entry.Log.Error("Failed to close stream", "error", err)
			}
			found = false
		}
		if !found {
			remoteIP, _, err := net.SplitHostPort(c.Request().RemoteAddr)
			if err != nil {
				remoteIP = c.Request().RemoteAddr
			}

			// Entry was not found
			entry = streamcache.NewEntry(remoteIP, presetEncoded, uuid)
			entry.Log.Info("Create stream")

			presetDecoded, err := preset.FromParam(presetEncoded)
			if err != nil {
				return apis.NewBadRequestError("", nil)
			}
			if len(presetDecoded.Tracks) == 0 {
				return apis.NewNotFoundError("", nil)
			}

			// Set up stream
			if entry.Streams, err = stream.New(dataFs, app.Dao(), presetDecoded); err != nil {
				if errors.Is(err, os.ErrNotExist) {
					// Invalid file ID returns 404
					return apis.NewNotFoundError("", nil)
				}
				// Other errors return 500
				panic(err)
			}
			if len(entry.Streams) == 0 {
				return apis.NewNotFoundError("", nil)
			}
			entry.Mix = entry.Streams.Mix()

			entry.Format = beep.Format{
				SampleRate:  44100,
				NumChannels: 2,
				Precision:   2,
			}

			if err := encode.VerifyFormat(entry.Format); err != nil {
				panic(err)
			}

			// Get current filetype encoder
			entry.Encoder, err = fileType.NewEncoder(entry.Buffer, entry.Format)
			if err != nil {
				panic(err)
			}

			cache.Add(uuid, entry)
		}

		c.Response().Header().Set("Accept-Ranges", "bytes")
		c.Response().Header().Set("Connection", "Keep-Alive")
		c.Response().Header().Set("Content-Type", fileType.ContentType())

		var firstByteIdx int
		if rangeHeader := c.Request().Header.Get("Range"); rangeHeader == "" {
			firstByteIdx = 0
		} else {
			unit, ranges, found := strings.Cut(rangeHeader, "=")
			// Error if no `=`, invalid unit, or multipart range
			if !found || unit != "bytes" || strings.ContainsRune(ranges, ',') {
				return apis.NewApiError(http.StatusRequestedRangeNotSatisfiable, "", nil)
			}

			first, last, found := strings.Cut(ranges, "-")
			// Error if `-` missing or end byte requested
			if !found || last != "" {
				return apis.NewApiError(http.StatusRequestedRangeNotSatisfiable, "", nil)
			}

			// Convert to int
			if firstByteIdx, err = strconv.Atoi(first); err != nil {
				return apis.NewApiError(http.StatusRequestedRangeNotSatisfiable, "", nil)
			}
		}

		var currChunkLength time.Duration
		// First chunks will be smaller to minimize delay
		switch entry.ChunkNum {
		case 0:
			currChunkLength = 15 * time.Second
		case 1:
			currChunkLength = time.Minute
		default:
			currChunkLength = chunkLength
		}

		// Ensure a single stream isn't fetched in parallel
		entry.Mu.Lock()
		defer entry.Mu.Unlock()

		// Mux streams to encoder
		if err = encode.Encode(c.Request().Context(), currChunkLength, entry); err != nil {
			panic(err)
		}

		if entry.TotalSize == 0 {
			// Set total length to ~24 hours
			// Actual length will vary when using VBR
			entry.TotalSize = entry.Buffer.Len() * int(24*time.Hour/currChunkLength)
		}

		c.Response().Header().Set("Content-Length", strconv.Itoa(entry.Buffer.Len()))
		lastByteIdx := firstByteIdx + entry.Buffer.Len() - 1
		c.Response().Header().Set(
			"Content-Range",
			fmt.Sprintf("bytes %d-%d/%d", firstByteIdx, lastByteIdx, entry.TotalSize),
		)

		// Write buffered stream data to client
		c.Response().WriteHeader(http.StatusPartialContent)
		n, err := io.Copy(c.Response(), entry.Buffer)
		if err != nil {
			if !errors.Is(err, syscall.EPIPE) && !errors.Is(err, syscall.ECONNRESET) {
				panic(err)
			}
		}
		entry.Transferred += uint64(n)

		entry.Buffer.Reset()
		entry.ChunkNum++
		return nil
	}
}
