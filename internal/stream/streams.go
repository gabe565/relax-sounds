package stream

import (
	"errors"
	"io/fs"

	"gabe565.com/relax-sounds/internal/config"
	"gabe565.com/relax-sounds/internal/preset"
	"github.com/gopxl/beep/v2"
)

//nolint:recvcheck
type Streams []Streamer

func (stream Streams) Close() error {
	errs := make([]error, 0, len(stream))
	for _, streamer := range stream {
		errs = append(errs, streamer.Close())
	}
	return errors.Join(errs...)
}

func (stream *Streams) Add(conf *config.Config, f fs.File, entry preset.Track) error {
	if entry.GetVolume() == 0 {
		return nil
	}

	streamer, err := NewStreamer(conf, f, entry)
	if err != nil {
		return err
	}

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

	if len(streams) == 1 {
		return streams[0]
	}
	return beep.Mix(streams...)
}
