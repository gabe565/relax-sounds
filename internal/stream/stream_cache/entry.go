package stream_cache

import (
	"bytes"
	"errors"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/gabe565/relax-sounds/internal/encoder"
	"github.com/gabe565/relax-sounds/internal/stream"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var activeStreamMetrics = promauto.NewGauge(prometheus.GaugeOpts{
	Namespace: "relax_sounds",
	Name:      "active_streams",
	Help:      "Active stream count",
})

var totalStreamMetrics = promauto.NewCounter(prometheus.CounterOpts{
	Namespace: "relax_sounds",
	Name:      "total_streams",
	Help:      "Total stream count",
})

type Entry struct {
	RemoteAddr string

	Preset  string
	Streams stream.Streams
	Mix     beep.Streamer
	Format  beep.Format

	Buffer  bytes.Buffer
	Encoder encoder.Encoder

	Mu        sync.Mutex
	ChunkNum  int
	TotalSize int
	Created   time.Time
	Accessed  time.Time
}

func NewEntry(remoteAddr, preset string) *Entry {
	entry := &Entry{
		RemoteAddr: remoteAddr,
		Preset:     preset,
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

	return errors.Join(
		e.Streams.Close(),
		e.Encoder.Close(),
	)
}
