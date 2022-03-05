package stream

import (
	"github.com/faiface/beep"
)

type Streamer struct {
	Streamer beep.Streamer
	Closer   beep.StreamCloser
	Format   beep.Format
}

func (s Streamer) Close() error {
	if s.Closer != nil {
		return s.Closer.Close()
	}
	return nil
}
