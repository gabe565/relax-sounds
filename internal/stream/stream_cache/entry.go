package stream_cache

import (
	"bytes"
	"errors"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/gabe565/relax-sounds/internal/encoder"
	"github.com/gabe565/relax-sounds/internal/stream"
)

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
	return entry
}

func (e *Entry) Close() error {
	e.Mu.Lock()
	defer e.Mu.Unlock()

	return errors.Join(
		e.Streams.Close(),
		e.Encoder.Close(),
	)
}
