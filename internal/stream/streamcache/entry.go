package streamcache

import (
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gabe565/relax-sounds/internal/encoder"
	"github.com/gabe565/relax-sounds/internal/stream"
	"github.com/gopxl/beep/v2"
	"github.com/labstack/echo/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

//nolint:gochecknoglobals
var (
	activeStreamMetrics = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "relax_sounds",
		Name:      "streams_active",
		Help:      "Active stream count",
	})

	totalStreamMetrics = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "relax_sounds",
		Name:      "streams_total",
		Help:      "Total stream count",
	})
)

type Entry struct {
	RemoteAddr string
	Log        *slog.Logger

	Preset  string
	Streams stream.Streams
	Mix     beep.Streamer
	Format  beep.Format

	Buffer  *Buffer
	Encoder encoder.Encoder

	Mu          sync.Mutex
	ChunkNum    int
	TotalSize   int
	Transferred uint64
	Created     time.Time
	Accessed    time.Time
}

func NewEntry(c echo.Context, preset, uuid string) *Entry {
	remoteIP := c.RealIP()
	entry := &Entry{
		RemoteAddr: remoteIP,
		Log: slog.With(
			"userIp", remoteIP,
			"userAgent", c.Request().UserAgent(),
			"url", c.Request().URL.String(),
			"id", uuid,
		),
		Preset:  preset,
		Buffer:  NewBuffer(3 * 1024 * 1024),
		Created: time.Now(),
	}
	entry.Accessed = entry.Created

	activeStreamMetrics.Inc()
	totalStreamMetrics.Inc()

	return entry
}

func (e *Entry) Close() error {
	e.Mu.Lock()
	defer e.Mu.Unlock()

	e.Log.Info("Close stream",
		"accessed", e.Accessed,
		"age", time.Since(e.Created).Round(100*time.Millisecond).String(),
		"transferred", humanize.IBytes(e.Transferred),
	)

	activeStreamMetrics.Dec()

	defer func() {
		e.Buffer = nil
	}()

	return errors.Join(
		e.Streams.Close(),
		e.Encoder.Close(),
	)
}
