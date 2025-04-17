package mix

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"gabe565.com/relax-sounds/internal/config"
	"gabe565.com/relax-sounds/internal/handlers/mix/cache"
	"gabe565.com/relax-sounds/internal/handlers/mix/encoder"
	"gabe565.com/relax-sounds/internal/handlers/mix/preset"
	"gabe565.com/relax-sounds/internal/handlers/mix/stream"
	"github.com/gopxl/beep/v2"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/valkey-io/valkey-go"
)

func NewMix(conf *config.Config) (*Mix, error) {
	var v valkey.Client
	if conf.ValkeyEnabled {
		valkeyAddr := net.JoinHostPort(conf.ValkeyHost, strconv.Itoa(int(conf.ValkeyPort)))
		var err error
		v, err = valkey.NewClient(valkey.ClientOption{
			InitAddress:           []string{valkeyAddr},
			Password:              conf.ValkeyPassword,
			SelectDB:              conf.ValkeyDB,
			DisableCache:          true,
			DisableAutoPipelining: true,
		})
		if err != nil {
			return nil, err
		}
	}

	return &Mix{
		conf:   conf,
		cache:  cache.New(conf),
		valkey: v,
	}, nil
}

type Mix struct {
	conf   *config.Config
	cache  *cache.Cache
	valkey valkey.Client
}

func (m *Mix) RegisterRoutes(e *core.ServeEvent) {
	g := e.Router.Group("/api/mix")
	g.Bind(apis.SkipSuccessActivityLog())
	g.GET("/{uuid}/{query}", m.Mix())
	g.DELETE("/{uuid}", m.Stop())
}

func (m *Mix) Mix() func(*core.RequestEvent) error { //nolint:gocyclo,gocognit,cyclop,funlen
	return func(e *core.RequestEvent) error {
		var err error

		query := e.Request.PathValue("query")

		// Preset parameter
		preStr, formatStr, found := strings.Cut(query, ".")
		if !found {
			return apis.NewNotFoundError("Missing file format", nil)
		}

		// File type parameter
		format, err := encoder.FormatString(formatStr)
		if err != nil {
			return apis.NewNotFoundError("Invalid file format", nil)
		}

		uuid := e.Request.PathValue("uuid")

		entry, found := m.cache.Get(uuid)
		if found && entry.Preset == preStr {
			// Ensure a single stream isn't fetched in parallel
			entry.Mu.Lock()
			defer entry.Mu.Unlock()
		} else {
			// Entry was not found
			_, _ = m.cache.Delete(uuid)

			pre, err := preset.FromParam(preStr)
			switch {
			case err != nil:
				return apis.NewBadRequestError("Failed to decode preset", nil)
			case len(pre) == 0:
				return apis.NewNotFoundError("Minimum preset length is 1 sound", nil)
			case len(pre) > m.conf.MaxPresetLen:
				return apis.NewBadRequestError("Maximum preset length is "+strconv.Itoa(m.conf.MaxPresetLen)+" sounds", nil)
			}

			var success bool
			entry = cache.NewEntry(e, uuid, preStr)
			defer func() {
				if !success {
					_ = entry.Close()
				}
			}()

			// Ensure a single stream isn't fetched in parallel
			entry.Mu.Lock()
			defer entry.Mu.Unlock()

			entry.Log.Info("Create stream")

			// Set up stream
			if entry.Streams, err = stream.New(m.conf, pre); err != nil {
				if errors.Is(err, os.ErrNotExist) {
					// Invalid file ID returns 404
					return apis.NewNotFoundError("Sounds not found", nil)
				}
				return apis.NewInternalServerError("Failed to create stream", nil)
			}
			if len(entry.Streams) == 0 {
				return apis.NewNotFoundError("Minimum preset length is 1 sound", nil)
			}

			entry.Format = beep.Format{
				SampleRate:  44100,
				NumChannels: 2,
				Precision:   2,
			}

			// Get current encoder
			entry.Encoder, err = format.NewEncoder(m.conf, entry.Writer, entry.Format)
			if err != nil {
				return apis.NewInternalServerError("Failed to create encoder", nil)
			}

			_ = m.cache.Set(uuid, entry)
			success = true
		}

		e.Response.Header().Set("Accept-Ranges", "bytes")
		e.Response.Header().Set("Connection", "Keep-Alive")
		e.Response.Header().Set("Content-Type", format.ContentType())

		var hasRangeHeader bool
		var chunkStart, chunkEnd int
		if rangeHeader := e.Request.Header.Get("Range"); rangeHeader != "" {
			hasRangeHeader = true

			unit, ranges, found := strings.Cut(rangeHeader, "=")
			// Error if no `=`, invalid unit, or multipart range
			if !found || unit != "bytes" || strings.ContainsRune(ranges, ',') {
				return apis.NewApiError(http.StatusRequestedRangeNotSatisfiable, "Missing unit", nil)
			}

			first, last, found := strings.Cut(ranges, "-")
			// Error if `-` missing or end byte requested
			if !found || last != "" {
				return apis.NewApiError(http.StatusRequestedRangeNotSatisfiable, "Missing chunk end byte", nil)
			}

			// Convert to int
			if chunkStart, err = strconv.Atoi(first); err != nil {
				return apis.NewApiError(http.StatusRequestedRangeNotSatisfiable, "Failed to parse chunk start", nil)
			}

			// Error if too large
			if chunkStart >= int(m.conf.MixTotalSize) {
				return apis.NewApiError(http.StatusRequestedRangeNotSatisfiable, "Range too large", nil)
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
			if chunkStart == 0 {
				if entry.Writer.TotalWritten() != 0 {
					for _, s := range entry.Streams {
						_ = s.Closer.Seek(0)
					}
				}
			} else if m.valkey != nil {
				if _, err := entry.LoadPositions(e.Request.Context(), m.valkey); err != nil {
					slog.Warn("Failed to load positions", "error", err)
				}
			}

			e.Response.Header().Set("Content-Length", strconv.Itoa(chunkSize))
			entry.Writer.Limit = chunkSize

			entry.Writer.SetWriter(e.Response)
			defer entry.Writer.SetWriter(nil)

			// Mux streams to encoder
			if err := Encode(e.Request.Context(), entry); err != nil {
				switch {
				case errors.Is(err, context.Canceled):
				case errors.Is(err, io.ErrShortWrite):
				case errors.Is(err, syscall.EPIPE):
				case errors.Is(err, syscall.ECONNRESET):
				default:
					return apis.NewInternalServerError("Failed to encode", nil)
				}
			}

			if m.valkey != nil {
				if err := entry.StorePositions(e.Request.Context(), m.valkey); err != nil {
					slog.Warn("Failed to store positions", "error", err)
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

		found, err := m.cache.Delete(uuid)
		if err != nil {
			return apis.NewInternalServerError("Failed to close cache entry", nil)
		}

		if !found {
			return apis.NewNotFoundError("No active stream for UUID", nil)
		}

		e.Response.WriteHeader(http.StatusNoContent)
		return nil
	}
}
