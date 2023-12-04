package stream

import (
	"errors"
	"io/fs"
	"sync"

	"github.com/gabe565/relax-sounds/internal/preset"
	"github.com/gopxl/beep"
)

type Streams []Streamer

func (stream Streams) Close() error {
	errs := make([]error, 0, len(stream))
	for _, streamer := range stream {
		errs = append(errs, streamer.Close())
	}
	return errors.Join(errs...)
}

func (stream *Streams) Add(f fs.File, entry preset.Track, mu *sync.Mutex) error {
	if entry.Volume == 0 {
		return nil
	}

	streamer, err := NewStreamer(f, entry)
	if err != nil {
		return err
	}

	mu.Lock()
	defer mu.Unlock()
	*stream = append(*stream, streamer)
	return nil
}

func (stream Streams) Mix() beep.Streamer {
	streams := make([]beep.Streamer, 0, len(stream))
	for _, streamer := range stream {
		switch {
		case streamer.Streamer != nil:
			streams = append(streams, streamer.Streamer)
		case streamer.Closer != nil:
			streams = append(streams, streamer.Closer)
		}
	}
	return beep.Mix(streams...)
}
