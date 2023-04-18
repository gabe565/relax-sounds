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
	f := File{File: rawFile}

	closer, format, err := f.Decode()
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

	if format.NumChannels < 2 {
		// Fix mono streams playing at 2x speed
		beepStreamer = beep.ResampleRatio(3, 0.5, beepStreamer)
	}
	streamer.Streamer = beepStreamer

	return streamer, nil
}
