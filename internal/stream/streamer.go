package stream

import (
	"io/fs"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/gabe565/relax-sounds/internal/preset"
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

func NewStreamer(rawFile fs.File, entry preset.Track) (streamer Streamer, err error) {
	closer, format, err := Decode(rawFile)
	if err != nil {
		return streamer, err
	}
	streamer.Closer = closer
	streamer.Format = format

	beepStreamer := beep.Loop(-1, closer)

	if entry.Volume != 1 {
		beepStreamer = &effects.Volume{
			Streamer: beepStreamer,
			Base:     10,
			Volume:   entry.Volume - 1,
		}
	}

	streamer.Streamer = beepStreamer
	return streamer, nil
}
