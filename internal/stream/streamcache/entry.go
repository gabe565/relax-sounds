package streamcache

import (
	"errors"
	"log/slog"
	"sync"
	"time"

	"gabe565.com/relax-sounds/internal/encoder"
	"gabe565.com/relax-sounds/internal/stream"
	"github.com/dustin/go-humanize"
	"github.com/gopxl/beep/v2"
	"github.com/pocketbase/pocketbase/core"
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
	Transferred uint64
	Created     time.Time
	Accessed    time.Time
}

func NewEntry(e *core.RequestEvent, preset, uuid string) *Entry {
	entry := &Entry{
		Log: slog.With(
			"userIp", e.RealIP(),
			"userAgent", e.Request.UserAgent(),
			"url", e.Request.URL.String(),
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
		"transferred", humanize.IBytes(e.Transferred),
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
