package stream

import (
	"io/fs"

	"gabe565.com/relax-sounds/internal/config"
	"gabe565.com/relax-sounds/internal/handlers/mix/preset"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/effects"
)

type Streamer struct {
	Streamer beep.Streamer
	Closer   beep.StreamSeekCloser
}

func (s Streamer) Close() error {
	if s.Closer != nil {
		return s.Closer.Close()
	}
	return nil
}

func NewStreamer(conf *config.Config, rawFile fs.File, entry preset.Track) (Streamer, error) {
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

	if rate := entry.GetRate(); rate != 1 || format.SampleRate != 44100 {
		rate *= float64(format.SampleRate) / 44100
		beepStreamer = beep.ResampleRatio(conf.ResampleQuality, rate, beepStreamer)
	}

	if volume := entry.GetVolume(); volume != 1 {
		beepStreamer = &effects.Volume{
			Streamer: beepStreamer,
			Base:     10,
			Volume:   volume - 1,
			Silent:   volume == 0,
		}
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
