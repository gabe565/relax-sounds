package stream

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/vorbis"
	"github.com/gabe565/relax-sounds/internal/preset"
	"io/fs"
	"sync"
)

type Streams []Streamer

func (stream Streams) Close() error {
	for _, streamer := range stream {
		_ = streamer.Close()
	}
	return nil
}

func (stream *Streams) Add(dataDir fs.FS, entry preset.Track, mu *sync.Mutex) error {
	infile, err := dataDir.Open(entry.Path())
	if err != nil {
		return err
	}

	closer, format, err := vorbis.Decode(infile)
	if err != nil {
		return err
	}

	streamer := beep.Loop(-1, closer)
	streamer = &effects.Volume{
		Streamer: streamer,
		Base:     4,
		Volume:   entry.Volume,
	}
	if format.NumChannels < 2 {
		// Fix mono streams playing at 2x speed
		streamer = beep.ResampleRatio(3, 0.5, streamer)
	}

	mu.Lock()
	defer mu.Unlock()
	*stream = append(*stream, Streamer{
		Streamer: streamer,
		Closer:   closer,
		Format:   format,
	})
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
