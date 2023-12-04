package stream

import (
	"io/fs"

	"github.com/gabe565/relax-sounds/internal/preset"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
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

	if volume := entry.GetVolume(); volume != 1 {
		beepStreamer = &effects.Volume{
			Streamer: beepStreamer,
			Base:     10,
			Volume:   volume - 1,
			Silent:   volume == 0,
		}
	}

	if rate := entry.GetRate(); rate != 1 {
		beepStreamer = beep.ResampleRatio(resampleQuality, rate, beepStreamer)
	}

	if pan := entry.GetPan(); pan != 0 {
		beepStreamer = &effects.Pan{
			Streamer: beepStreamer,
			Pan:      pan,
		}
	}

	streamer.Streamer = beepStreamer
	return streamer, nil
}
