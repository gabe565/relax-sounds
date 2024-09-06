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

	"github.com/gabe565/relax-sounds/internal/encoder/encode"
	"github.com/gabe565/relax-sounds/internal/encoder/filetype"
	"github.com/gabe565/relax-sounds/internal/preset"
	"github.com/gabe565/relax-sounds/internal/stream"
	"github.com/gabe565/relax-sounds/internal/stream/streamcache"
	"github.com/gopxl/beep/v2"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
)

func NewMix(app *pocketbase.PocketBase) *Mix {
	return &Mix{
		app:   app,
		cache: streamcache.New(),
	}
}

type Mix struct {
	app   *pocketbase.PocketBase
	cache *streamcache.Cache
}

func (m *Mix) RegisterRoutes(e *echo.Echo) {
	e.GET("/api/mix/:uuid/:query", m.Mix())
	e.DELETE("/api/mix/:uuid", m.Stop(), middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(1)))
}

func (m *Mix) Mix() echo.HandlerFunc {
	dataFs := os.DirFS(filepath.Join(m.app.DataDir(), "storage"))

	const (
		TotalSize = 1500 * 1024 * 1024
		ChunkSize = 2 * 1024 * 1024
	)

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
		entry, found := m.cache.Get(uuid)
		if found && entry.Preset != presetEncoded {
			// Same window is changing streams
			// Destroy old stream then recreate
			if err := entry.Close(); err != nil {
				entry.Log.Error("Failed to close stream", "error", err)
			}
			found = false
		}
		if !found {
			// Entry was not found
			entry = streamcache.NewEntry(c, presetEncoded, uuid)
			entry.Log.Info("Create stream")

			presetDecoded, err := preset.FromParam(presetEncoded)
			if err != nil {
				return apis.NewBadRequestError("", nil)
			}
			if len(presetDecoded.Tracks) == 0 {
				return apis.NewNotFoundError("", nil)
			}

			// Set up stream
			if entry.Streams, err = stream.New(dataFs, m.app.Dao(), presetDecoded); err != nil {
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

			// Get current filetype encoder
			entry.Encoder, err = fileType.NewEncoder(entry.Writer, entry.Format)
			if err != nil {
				panic(err)
			}

			m.cache.Add(uuid, entry)
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

		// Ensure a single stream isn't fetched in parallel
		entry.Mu.Lock()
		defer entry.Mu.Unlock()

		chunkSize := ChunkSize + entry.Writer.Buffered()
		c.Response().Header().Set("Content-Length", strconv.Itoa(chunkSize))
		c.Response().Header().Set("Content-Range", fmt.Sprintf(
			"bytes %d-%d/%d",
			firstByteIdx,
			firstByteIdx+chunkSize-1,
			TotalSize,
		))
		entry.Writer.SetWriter(c.Response())
		entry.Writer.Limit = chunkSize

		c.Response().WriteHeader(http.StatusPartialContent)
		defer func() {
			entry.Transferred += c.Response().Size
		}()

		// Mux streams to encoder
		if err = encode.Encode(c.Request().Context(), entry); err != nil {
			switch {
			case errors.Is(err, io.ErrShortWrite):
			case errors.Is(err, syscall.EPIPE):
			case errors.Is(err, syscall.ECONNRESET):
			default:
				return apis.NewApiError(http.StatusInternalServerError, "", err)
			}
		}

		return nil
	}
}

func (m *Mix) Stop() echo.HandlerFunc {
	return func(c echo.Context) error {
		uuid := c.PathParam("uuid")
		entry, found := m.cache.Get(uuid)
		if !found {
			return apis.NewNotFoundError("no active stream for uuid", nil)
		}

		m.cache.Delete(uuid)
		if err := entry.Close(); err != nil {
			return apis.NewApiError(http.StatusInternalServerError, "failed to close cache entry", nil)
		}

		c.Response().WriteHeader(http.StatusNoContent)
		return nil
	}
}
