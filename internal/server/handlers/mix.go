package handlers

import (
	"errors"
	"fmt"
	"github.com/faiface/beep"
	"github.com/gabe565/relax-sounds/internal/encoder/filetype"
	"github.com/gabe565/relax-sounds/internal/mix"
	"github.com/gabe565/relax-sounds/internal/preset"
	"github.com/gabe565/relax-sounds/internal/stream"
	"github.com/gabe565/relax-sounds/internal/stream/stream_cache"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func Mix() http.HandlerFunc {
	cache := stream_cache.New()

	return func(res http.ResponseWriter, req *http.Request) {
		var err error
		ctx := req.Context()

		uuid := chi.URLParam(req, "uuid")
		presetEncoded := chi.URLParam(req, "enc")
		preset := ctx.Value(preset.RequestKey).(preset.Preset)
		fileType := ctx.Value(filetype.RequestKey).(filetype.FileType)

		entry, found := cache.Get(uuid)
		if found && entry.Preset != presetEncoded {
			// Same window is changing streams
			// Destroy old stream then recreate
			log.WithField("id", uuid).Info("Close stream")
			if err := entry.Close(); err != nil {
				log.WithError(err).WithField("id", uuid).
					Error("Failed to close stream")
			}
			found = false
		}
		if !found {
			log.WithField("id", uuid).Info("Create stream")

			// Entry was not found
			entry = stream_cache.NewEntry(presetEncoded)

			// Set up stream
			if entry.Streams, err = stream.New(preset); err != nil {
				if errors.Is(err, os.ErrNotExist) {
					// Invalid file ID returns 404
					http.Error(res, http.StatusText(404), 404)
					return
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

			if err := mix.VerifyFormat(entry.Format); err != nil {
				panic(err)
			}

			// Get current filetype encoder
			entry.Encoder, err = fileType.NewEncoder(&entry.Buffer, entry.Format)
			if err != nil {
				panic(err)
			}

			cache.Add(uuid, entry)
		}

		res.Header().Set("Accept-Ranges", "bytes")
		res.Header().Set("Connection", "Keep-Alive")
		res.Header().Set("Content-Type", fileType.ContentType())

		var firstByteIdx int
		if rangeHeader := req.Header.Get("Range"); rangeHeader == "" {
			// Write 200 OK if no range requested
			res.WriteHeader(http.StatusOK)
			return
		} else {
			unit, ranges, found := strings.Cut(rangeHeader, "=")
			// Error if no `=`, invalid unit, or multipart range
			if !found || unit != "bytes" || strings.ContainsRune(ranges, ',') {
				http.Error(res, http.StatusText(http.StatusRequestedRangeNotSatisfiable), http.StatusRequestedRangeNotSatisfiable)
				return
			}

			first, last, found := strings.Cut(ranges, "-")
			// Error if `-` missing or end byte requested
			if !found || last != "" {
				http.Error(res, http.StatusText(http.StatusRequestedRangeNotSatisfiable), http.StatusRequestedRangeNotSatisfiable)
				return
			}

			// Convert to int
			if firstByteIdx, err = strconv.Atoi(first); err != nil {
				http.Error(res, http.StatusText(http.StatusRequestedRangeNotSatisfiable), http.StatusRequestedRangeNotSatisfiable)
				return
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
		if err = mix.Encode(ctx, chunkSize, entry); err != nil {
			panic(err)
		}

		if entry.TotalSize == 0 {
			// Set total length to ~16 hours
			// Actual length will vary when using VBR
			entry.TotalSize = entry.Buffer.Len() * int(16*time.Hour/chunkSize)
		}

		res.Header().Set("Content-Length", strconv.Itoa(entry.Buffer.Len()))
		lastByteIdx := firstByteIdx + entry.Buffer.Len() - 1
		res.Header().Set(
			"Content-Range",
			fmt.Sprintf("bytes %d-%d/%d", firstByteIdx, lastByteIdx, entry.TotalSize),
		)

		// Write buffered stream data to client
		res.WriteHeader(http.StatusPartialContent)
		if _, err := io.Copy(res, &entry.Buffer); err != nil {
			if !errors.Is(err, syscall.EPIPE) && !errors.Is(err, syscall.ECONNRESET) {
				panic(err)
			}
		}

		entry.Buffer.Reset()
		entry.ChunkNum += 1
	}
}
