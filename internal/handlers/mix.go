package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

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
	g := e.Router.Group("/api/mix")
	g.Bind(apis.SkipSuccessActivityLog())
	g.GET("/{uuid}/{query}", m.Mix())
	g.DELETE("/{uuid}", m.Stop())
}

func (m *Mix) Mix() func(*core.RequestEvent) error { //nolint:gocyclo
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
		if found && entry.Preset == presetEncoded {
			// Ensure a single stream isn't fetched in parallel
			entry.Mu.Lock()
			defer entry.Mu.Unlock()
		} else {
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

			// Ensure a single stream isn't fetched in parallel
			entry.Mu.Lock()
			defer entry.Mu.Unlock()

			if err := m.cache.Set(uuid, entry); err != nil {
				entry.Log.Error("Failed to close stream", "error", err)
			}

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
		}

		e.Response.Header().Set("Accept-Ranges", "bytes")
		e.Response.Header().Set("Connection", "Keep-Alive")
		e.Response.Header().Set("Content-Type", fileType.ContentType())

		var hasRangeHeader bool
		var chunkStart, chunkEnd int
		if rangeHeader := e.Request.Header.Get("Range"); rangeHeader != "" {
			hasRangeHeader = true

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
			if chunkStart, err = strconv.Atoi(first); err != nil {
				return apis.NewApiError(http.StatusRequestedRangeNotSatisfiable, "", nil)
			}
		}

		chunkSize := int(m.conf.MixChunkSize) + entry.Writer.Buffered()
		if hasRangeHeader {
			chunkEnd = chunkStart + chunkSize - 1
			if chunkEnd >= int(m.conf.MixTotalSize) {
				chunkEnd = int(m.conf.MixTotalSize) - 1
				chunkSize = chunkEnd - chunkStart + 1
			}
		}
		e.Response.Header().Set("Content-Range",
			"bytes "+strconv.Itoa(chunkStart)+"-"+strconv.Itoa(chunkEnd)+"/"+strconv.Itoa(int(m.conf.MixTotalSize)),
		)

		e.Response.WriteHeader(http.StatusPartialContent)

		if hasRangeHeader && e.Request.Method == http.MethodGet {
			if chunkStart == 0 && entry.Writer.TotalWritten() != 0 {
				for _, s := range entry.Streams {
					_ = s.Closer.Seek(0)
				}
			}

			e.Response.Header().Set("Content-Length", strconv.Itoa(chunkSize))
			entry.Writer.Limit = chunkSize

			entry.Writer.SetWriter(e.Response)
			defer entry.Writer.SetWriter(nil)

			// Mux streams to encoder
			if err := encode.Encode(e.Request.Context(), entry); err != nil {
				switch {
				case errors.Is(err, context.Canceled):
				case errors.Is(err, io.ErrShortWrite):
				case errors.Is(err, syscall.EPIPE):
				case errors.Is(err, syscall.ECONNRESET):
				default:
					return apis.NewApiError(http.StatusInternalServerError, "", err)
				}
			}
		}

		entry.Accessed = time.Now()
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
