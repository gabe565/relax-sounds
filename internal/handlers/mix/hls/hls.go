package hls

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gabe565.com/relax-sounds/internal/config"
	"gabe565.com/relax-sounds/internal/handlers/mix/preset"
	"gabe565.com/relax-sounds/internal/handlers/mix/stream"
	"github.com/jellydator/ttlcache/v3"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func New(conf *config.Config) *HLS {
	return &HLS{
		conf:  conf,
		cache: newCache(conf),
	}
}

type HLS struct {
	conf  *config.Config
	cache *Cache
}

func (h *HLS) RegisterRoutes(e *core.ServeEvent) {
	g := e.Router.Group("/api/mix")
	g.Bind(apis.SkipSuccessActivityLog())
	g.GET("/{uuid}/{query}", h.Manifest())
	g.GET("/{uuid}/{preset}/{seq}", h.Segment())
	g.DELETE("/{uuid}", h.Stop())
}

// getOrCreateEntry returns the cached entry for (uuid, preset), creating and
// starting a new producer if none exists or the preset changed.
func (h *HLS) getOrCreateEntry(e *core.RequestEvent, uuid, presetStr string) (*Entry, error) {
	if cacheEntry := h.cache.Get(uuid); cacheEntry != nil && cacheEntry.Value().Preset == presetStr {
		return cacheEntry.Value(), nil
	}

	pre, err := preset.FromParam(presetStr)
	switch {
	case err != nil:
		return nil, apis.NewBadRequestError("Failed to decode preset", nil)
	case len(pre) == 0:
		return nil, apis.NewNotFoundError("Minimum preset length is 1 sound", nil)
	case len(pre) > h.conf.MaxPresetLen:
		return nil, apis.NewBadRequestError(
			"Maximum preset length is "+strconv.Itoa(h.conf.MaxPresetLen)+" sounds", nil,
		)
	}

	entry := NewEntry(e, uuid, presetStr)
	entry.Streams, err = stream.New(h.conf, pre)
	if err != nil {
		_ = entry.Close()
		if errors.Is(err, os.ErrNotExist) {
			return nil, apis.NewNotFoundError("Sounds not found", nil)
		}
		return nil, apis.NewInternalServerError("Failed to create stream", nil)
	}
	if len(entry.Streams) == 0 {
		_ = entry.Close()
		return nil, apis.NewNotFoundError("Minimum preset length is 1 sound", nil)
	}

	h.cache.Set(uuid, entry, ttlcache.DefaultTTL)
	go entry.Produce(h.conf)
	return entry, nil
}

// Manifest serves the HLS live playlist for a UUID/preset.
func (h *HLS) Manifest() func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		uuid := e.Request.PathValue("uuid")
		query := e.Request.PathValue("query")
		presetStr, ext, found := strings.Cut(query, ".")
		if !found || ext != "m3u8" {
			return apis.NewNotFoundError("Invalid manifest path", nil)
		}

		entry, err := h.getOrCreateEntry(e, uuid, presetStr)
		if err != nil {
			return err
		}

		// Block until the producer has accumulated enough segments for the
		// player to start playback without immediate underrun. Bounded by the
		// request's context.
		if err := entry.WaitReady(e.Request.Context()); err != nil {
			if errors.Is(err, ErrClosed) {
				return apis.NewInternalServerError("Stream closed before ready", nil)
			}
			return err
		}

		body := renderPlaylist(entry, presetStr)

		e.Response.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
		e.Response.Header().Set("Cache-Control", "no-store")
		e.Response.WriteHeader(http.StatusOK)
		_, _ = e.Response.Write([]byte(body))
		return nil
	}
}

// Segment serves an individual MP3 segment.
func (h *HLS) Segment() func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		uuid := e.Request.PathValue("uuid")
		presetStr := e.Request.PathValue("preset")
		seqWithExt := e.Request.PathValue("seq")
		seqStr, ext, found := strings.Cut(seqWithExt, ".")
		if !found || ext != "mp3" {
			return apis.NewNotFoundError("Invalid segment path", nil)
		}

		seq, err := strconv.ParseUint(seqStr, 10, 64)
		if err != nil {
			return apis.NewBadRequestError("Invalid segment sequence", nil)
		}

		cacheEntry := h.cache.Get(uuid)
		if cacheEntry == nil || cacheEntry.Value().Preset != presetStr {
			return apis.NewNotFoundError("Stream not found", nil)
		}
		seg := cacheEntry.Value().GetSegment(seq)
		if seg == nil {
			return apis.NewNotFoundError("Segment not available", nil)
		}

		e.Response.Header().Set("Content-Type", "audio/mpeg")
		e.Response.Header().Set("Content-Length", strconv.Itoa(len(seg.Bytes)))
		e.Response.Header().Set("Cache-Control", "no-store")
		e.Response.WriteHeader(http.StatusOK)
		n, _ := e.Response.Write(seg.Bytes)
		cacheEntry.Value().AddTransferred(int64(n))
		return nil
	}
}

// Stop tears down the stream for a UUID.
func (h *HLS) Stop() func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		uuid := e.Request.PathValue("uuid")
		h.cache.Delete(uuid)
		e.Response.WriteHeader(http.StatusNoContent)
		return nil
	}
}

// renderPlaylist builds a live HLS playlist for the given entry.
func renderPlaylist(entry *Entry, presetStr string) string {
	var b strings.Builder
	b.Grow(84 + ManifestWindow*(30+len(presetStr)))
	b.WriteString("#EXTM3U\n")
	b.WriteString("#EXT-X-VERSION:3\n")
	targetSecs := int(SegmentDuration().Round(time.Second).Seconds())
	fmt.Fprintf(&b, "#EXT-X-TARGETDURATION:%d\n", targetSecs)

	for i, s := range entry.buf.Last(ManifestWindow) {
		if i == 0 {
			fmt.Fprintf(&b, "#EXT-X-MEDIA-SEQUENCE:%d\n", s.Seq)
		}
		fmt.Fprintf(&b, "#EXTINF:%.3f,\n", s.Duration.Seconds())
		fmt.Fprintf(&b, "%s/%d.mp3\n", presetStr, s.Seq)
	}
	return b.String()
}
