package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"syscall"

	"gabe565.com/relax-sounds/internal/config"
	"gabe565.com/relax-sounds/internal/encoder/encode"
	"gabe565.com/relax-sounds/internal/encoder/filetype"
	"gabe565.com/relax-sounds/internal/preset"
	"gabe565.com/relax-sounds/internal/stream"
	"gabe565.com/relax-sounds/internal/stream/streamcache"
	"github.com/gopxl/beep/v2"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func NewMix(conf *config.Config) *Mix {
	return &Mix{
		conf:  conf,
		cache: streamcache.New(conf),
	}
}

type Mix struct {
	conf  *config.Config
	cache *streamcache.Cache
}

func (m *Mix) RegisterRoutes(e *core.ServeEvent) {
	e.Router.GET("/api/mix/{uuid}/{query}", m.Mix())
	e.Router.DELETE("/api/mix/{uuid}", m.Stop())
}

func (m *Mix) Mix() func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		var err error

		query := e.Request.PathValue("query")

		// Preset parameter
		presetEncoded, fileTypeStr, found := strings.Cut(query, ".")
		if !found {
			return apis.NewNotFoundError("", nil)
		}

		// File type parameter
		fileType, err := filetype.FileTypeString(fileTypeStr)
		if err != nil {
			return apis.NewNotFoundError("", nil)
		}

		uuid := e.Request.PathValue("uuid")
		entry, found := m.cache.Get(uuid)
		if !found || entry.Preset != presetEncoded {
			// Entry was not found
			presetDecoded, err := preset.FromParam(presetEncoded)
			switch {
			case err != nil:
				return apis.NewBadRequestError("", nil)
			case len(presetDecoded) == 0:
				return apis.NewNotFoundError("", nil)
			case len(presetDecoded) > m.conf.MaxPresetLen:
				return apis.NewBadRequestError("Maximum preset length is "+strconv.Itoa(m.conf.MaxPresetLen)+" sounds.", nil)
			}

			entry = streamcache.NewEntry(e, presetEncoded, uuid)
			entry.Log.Info("Create stream")

			// Set up stream
			if entry.Streams, err = stream.New(m.conf, presetDecoded); err != nil {
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

			entry.Format = beep.Format{
				SampleRate:  44100,
				NumChannels: 2,
				Precision:   2,
			}

			// Get current filetype encoder
			entry.Encoder, err = fileType.NewEncoder(m.conf, entry.Writer, entry.Format)
			if err != nil {
				panic(err)
			}

			if err := m.cache.Set(uuid, entry); err != nil {
				entry.Log.Error("Failed to close stream", "error", err)
			}
		}

		e.Response.Header().Set("Accept-Ranges", "bytes")
		e.Response.Header().Set("Connection", "Keep-Alive")
		e.Response.Header().Set("Content-Type", fileType.ContentType())

		var firstByteIdx int
		if rangeHeader := e.Request.Header.Get("Range"); rangeHeader != "" {
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

		chunkSize := int(m.conf.MixChunkSize) + entry.Writer.Buffered()
		e.Response.Header().Set("Content-Length", strconv.Itoa(chunkSize))
		e.Response.Header().Set("Content-Range", fmt.Sprintf(
			"bytes %d-%d/%d",
			firstByteIdx,
			firstByteIdx+chunkSize-1,
			int(m.conf.MixTotalSize),
		))
		entry.Writer.SetWriter(e.Response)
		defer entry.Writer.SetWriter(nil)
		entry.Writer.Limit = chunkSize

		e.Response.WriteHeader(http.StatusPartialContent)

		// Mux streams to encoder
		if err := encode.Encode(e.Request.Context(), entry); err != nil {
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
		if !m.cache.Has(uuid) {
			return apis.NewNotFoundError("no active stream for uuid", nil)
		}

		if err := m.cache.Delete(uuid); err != nil {
			return apis.NewApiError(http.StatusInternalServerError, "failed to close cache entry", nil)
		}

		e.Response.WriteHeader(http.StatusNoContent)
		return nil
	}
}
