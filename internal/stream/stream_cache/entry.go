package stream_cache

import (
	"bytes"
	"errors"
	"github.com/faiface/beep"
	"github.com/gabe565/relax-sounds/internal/encoder"
	"github.com/gabe565/relax-sounds/internal/stream"
	"sync"
	"time"
)

type Entry struct {
	Preset  string
	Streams stream.Streams
	Mix     beep.Streamer
	Format  beep.Format

	Buffer  bytes.Buffer
	Encoder encoder.Encoder

	Mu        sync.Mutex
	ChunkNum  int
	TotalSize int
	Accessed  time.Time
}

func NewEntry(preset string) *Entry {
	return &Entry{
		Preset: preset,
	}
}

func (e *Entry) Close() error {
	return errors.Join(
		e.Streams.Close(),
		e.Encoder.Close(),
	)
}
