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
	Log *slog.Logger

	Preset  string
	Streams stream.Streams
	Mix     beep.Streamer
	Format  beep.Format

	Writer  *ProxyWriter
	Encoder encoder.Encoder

	Mu          sync.Mutex
	Transferred int64
	Created     time.Time
	Accessed    time.Time
}

func NewEntry(c echo.Context, preset, uuid string) *Entry {
	entry := &Entry{
		Log: slog.With(
			"userIp", c.RealIP(),
			"userAgent", c.Request().UserAgent(),
			"url", c.Request().URL.String(),
			"id", uuid,
		),
		Preset:  preset,
		Writer:  NewProxyWriter(),
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
		"transferred", humanize.IBytes(uint64(e.Transferred)),
	)

	activeStreamMetrics.Dec()

	e.Writer.Close()
	defer func() {
		e.Writer = nil
	}()

	return errors.Join(
		e.Streams.Close(),
		e.Encoder.Close(),
	)
}
