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

	"gabe565.com/relax-sounds/internal/encoder/encode"
	"gabe565.com/relax-sounds/internal/encoder/filetype"
	"gabe565.com/relax-sounds/internal/preset"
	"gabe565.com/relax-sounds/internal/stream"
	"gabe565.com/relax-sounds/internal/stream/streamcache"
	"github.com/gopxl/beep/v2"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
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

func (m *Mix) RegisterRoutes(e *core.ServeEvent) {
	e.Router.GET("/api/mix/{uuid}/{query}", m.Mix())
	e.Router.DELETE("/api/mix/{uuid}", m.Stop())
}

func (m *Mix) Mix() func(*core.RequestEvent) error {
	dataFs := os.DirFS(filepath.Join(m.app.DataDir(), "storage"))

	const (
		TotalSize = 1500 * 1024 * 1024
		ChunkSize = 2 * 1024 * 1024
	)

	return func(e *core.RequestEvent) error {
		var err error

		query := e.Request.PathValue("query")

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

		uuid := e.Request.PathValue("uuid")
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
			entry = streamcache.NewEntry(e, presetEncoded, uuid)
			entry.Log.Info("Create stream")

			presetDecoded, err := preset.FromParam(presetEncoded)
			if err != nil {
				return apis.NewBadRequestError("", nil)
			}
			if len(presetDecoded.Tracks) == 0 {
				return apis.NewNotFoundError("", nil)
			}

			// Set up stream
			if entry.Streams, err = stream.New(dataFs, m.app, presetDecoded); err != nil {
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

		e.Response.Header().Set("Accept-Ranges", "bytes")
		e.Response.Header().Set("Connection", "Keep-Alive")
		e.Response.Header().Set("Content-Type", fileType.ContentType())

		var firstByteIdx int
		if rangeHeader := e.Request.Header.Get("Range"); rangeHeader == "" {
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
		e.Response.Header().Set("Content-Length", strconv.Itoa(chunkSize))
		e.Response.Header().Set("Content-Range", fmt.Sprintf(
			"bytes %d-%d/%d",
			firstByteIdx,
			firstByteIdx+chunkSize-1,
			TotalSize,
		))
		entry.Writer.SetWriter(e.Response)
		defer entry.Writer.SetWriter(nil)
		entry.Writer.Limit = chunkSize

		e.Response.WriteHeader(http.StatusPartialContent)

		// Mux streams to encoder
		n, err := encode.Encode(e.Request.Context(), entry)
		entry.Transferred += uint64(n) //nolint:gosec
		if err != nil {
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

func (m *Mix) Stop() func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		uuid := e.Request.PathValue("uuid")
		entry, found := m.cache.Get(uuid)
		if !found {
			return apis.NewNotFoundError("no active stream for uuid", nil)
		}

		m.cache.Delete(uuid)
		if err := entry.Close(); err != nil {
			return apis.NewApiError(http.StatusInternalServerError, "failed to close cache entry", nil)
		}

		e.Response.WriteHeader(http.StatusNoContent)
		return nil
	}
}
