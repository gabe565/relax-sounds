package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/faiface/beep"
	"github.com/gabe565/relax-sounds/internal/encoder/encode"
	"github.com/gabe565/relax-sounds/internal/encoder/filetype"
	"github.com/gabe565/relax-sounds/internal/preset"
	"github.com/gabe565/relax-sounds/internal/stream"
	"github.com/gabe565/relax-sounds/internal/stream/stream_cache"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	log "github.com/sirupsen/logrus"
)

func Mix(app core.App) echo.HandlerFunc {
	cache := stream_cache.New()
	dataFs := os.DirFS(filepath.Join(app.DataDir(), "storage"))

	return func(c echo.Context) error {
		var err error

		query := c.PathParam("query")

		// Preset parameter
		presetEncoded, fileTypeStr, found := strings.Cut(query, ".")
		if !found {
			return apis.NewNotFoundError("", nil)
		}

		preset, err := preset.FromParam(presetEncoded)
		if err != nil {
			return apis.NewBadRequestError("", nil)
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
			log.WithFields(log.Fields{
				"ip":       entry.RemoteAddr,
				"id":       uuid,
				"accessed": entry.Accessed,
				"age":      time.Since(entry.Created).Truncate(time.Millisecond).String(),
			}).Info("Close stream")
			if err := entry.Close(); err != nil {
				log.WithError(err).WithField("id", uuid).
					Error("Failed to close stream")
			}
			found = false
		}
		if !found {
			log.WithField("id", uuid).Info("Create stream")

			remoteAddr, _, _ := strings.Cut(c.Request().RemoteAddr, ":")

			// Entry was not found
			entry = stream_cache.NewEntry(remoteAddr, presetEncoded)

			// Set up stream
			if entry.Streams, err = stream.New(dataFs, app.Dao(), preset); err != nil {
				if errors.Is(err, os.ErrNotExist) {
					// Invalid file ID returns 404
					return apis.NewNotFoundError("", nil)
				}
				// Other errors return 500
				panic(err)
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
			entry.Encoder, err = fileType.NewEncoder(&entry.Buffer, entry.Format)
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

		chunkSize := 2 * time.Minute
		// First chunks will be smaller to minimize delay
		switch entry.ChunkNum {
		case 0:
			chunkSize = 15 * time.Second
		case 1:
			chunkSize = time.Minute
		}

		// Ensure a single stream isn't fetched in parallel
		entry.Mu.Lock()
		defer entry.Mu.Unlock()

		// Mux streams to encoder
		if err = encode.Encode(c.Request().Context(), chunkSize, entry); err != nil {
			panic(err)
		}

		if entry.TotalSize == 0 {
			// Set total length to ~16 hours
			// Actual length will vary when using VBR
			entry.TotalSize = entry.Buffer.Len() * int(16*time.Hour/chunkSize)
		}

		c.Response().Header().Set("Content-Length", strconv.Itoa(entry.Buffer.Len()))
		lastByteIdx := firstByteIdx + entry.Buffer.Len() - 1
		c.Response().Header().Set(
			"Content-Range",
			fmt.Sprintf("bytes %d-%d/%d", firstByteIdx, lastByteIdx, entry.TotalSize),
		)

		// Write buffered stream data to client
		c.Response().WriteHeader(http.StatusPartialContent)
		if _, err := io.Copy(c.Response(), &entry.Buffer); err != nil {
			if !errors.Is(err, syscall.EPIPE) && !errors.Is(err, syscall.ECONNRESET) {
				panic(err)
			}
		}

		entry.Buffer.Reset()
		entry.ChunkNum += 1
		return nil
	}
}
