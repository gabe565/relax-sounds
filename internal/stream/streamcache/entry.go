package streamcache

import (
	"bytes"
	"errors"
	"sync"
	"time"

	"github.com/gabe565/relax-sounds/internal/encoder"
	"github.com/gabe565/relax-sounds/internal/stream"
	"github.com/gopxl/beep"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	Log        zerolog.Logger

	Preset  string
	Streams stream.Streams
	Mix     beep.Streamer
	Format  beep.Format

	Buffer  *bytes.Buffer
	Encoder encoder.Encoder

	Mu        sync.Mutex
	ChunkNum  int
	TotalSize int
	Created   time.Time
	Accessed  time.Time
}

func NewEntry(remoteAddr, preset, uuid string) *Entry {
	entry := &Entry{
		RemoteAddr: remoteAddr,
		Log:        log.With().Str("ip", remoteAddr).Str("id", uuid).Logger(),
		Preset:     preset,
		Buffer:     bytes.NewBuffer(make([]byte, 0, 3*1024*1024)),
		Created:    time.Now(),
	}
	entry.Accessed = entry.Created

	activeStreamMetrics.Inc()
	totalStreamMetrics.Inc()

	return entry
}

func (e *Entry) Close() error {
	e.Mu.Lock()
	defer e.Mu.Unlock()

	activeStreamMetrics.Dec()

	defer func() {
		e.Buffer = nil
	}()

	return errors.Join(
		e.Streams.Close(),
		e.Encoder.Close(),
	)
}
