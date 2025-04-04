package streamcache

import (
	"errors"
	"log/slog"
	"sync"
	"time"

	"gabe565.com/relax-sounds/internal/config"
	"gabe565.com/relax-sounds/internal/encoder"
	"gabe565.com/relax-sounds/internal/stream"
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
	Format  beep.Format

	Writer  *ProxyWriter
	Encoder encoder.Encoder

	Mu       sync.Mutex
	Created  time.Time
	Accessed time.Time
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
		"transferred", config.Bytes(e.Writer.TotalWritten()).String(),
	)

	activeStreamMetrics.Dec()

	if e.Writer != nil {
		e.Writer.Close()
	}
	defer func() {
		e.Writer = nil
	}()

	errs := make([]error, 0, 2)
	errs = append(errs, e.Streams.Close())
	if e.Encoder != nil {
		errs = append(errs, e.Encoder.Close())
	}
	return errors.Join(errs...)
}
