package streamcache

import (
	"bytes"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/gabe565/relax-sounds/internal/encoder"
	"github.com/gabe565/relax-sounds/internal/stream"
	"github.com/gopxl/beep/v2"
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

	Buffer  *bytes.Buffer
	Encoder encoder.Encoder

	Mu          sync.Mutex
	ChunkNum    int
	TotalSize   int
	Transferred uint64
	Created     time.Time
	Accessed    time.Time
}

func NewEntry(remoteAddr, preset, uuid string) *Entry {
	entry := &Entry{
		RemoteAddr: remoteAddr,
		Log:        slog.With("userIp", remoteAddr, "id", uuid),
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
