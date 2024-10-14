package stream

import (
	"io/fs"

	"gabe565.com/relax-sounds/internal/preset"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
)

type Streamer struct {
	Streamer beep.Streamer
	Closer   beep.StreamCloser
}

func (s Streamer) Close() error {
	if s.Closer != nil {
		return s.Closer.Close()
	}
	return nil
}

func NewStreamer(rawFile fs.File, entry preset.Track) (Streamer, error) {
	var streamer Streamer

	closer, format, err := Decode(rawFile)
	if err != nil {
		return streamer, err
	}
	streamer.Closer = closer

	beepStreamer, err := beep.Loop2(closer)
	if err != nil {
		return streamer, err
	}

	if format.SampleRate != 44100 {
		beepStreamer = beep.Resample(3, format.SampleRate, 44100, beepStreamer)
	}

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
