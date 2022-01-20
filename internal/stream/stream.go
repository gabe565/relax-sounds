package stream

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/vorbis"
	"github.com/gabe565/relax-sounds/internal/preset"
	"golang.org/x/sync/errgroup"
	"io/fs"
	"sync"
)

type Stream struct {
	closers   []beep.StreamSeekCloser
	streamers []beep.Streamer
	Formats   []beep.Format
}

func (stream *Stream) Close() error {
	for _, streamer := range stream.closers {
		_ = streamer.Close()
	}
	return nil
}

func (stream *Stream) Add(dataDir fs.FS, entry preset.Track, mu *sync.Mutex) error {
	infile, err := dataDir.Open(entry.Path())
	if err != nil {
		return err
	}

	streamCloser, format, err := vorbis.Decode(infile)
	if err != nil {
		return err
	}

	streamer := beep.Loop(-1, streamCloser)
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
	stream.closers = append(stream.closers, streamCloser)
	stream.streamers = append(stream.streamers, streamer)
	stream.Formats = append(stream.Formats, format)
	return nil
}

func (stream *Stream) Mix() beep.Streamer {
	return beep.Mix(stream.streamers...)
}

func New(p preset.Preset) (stream Stream, err error) {
	s := Stream{
		closers:   make([]beep.StreamSeekCloser, 0, len(p.Tracks)),
		streamers: make([]beep.Streamer, 0, len(p.Tracks)),
		Formats:   make([]beep.Format, 0, len(p.Tracks)),
	}

	var mu sync.Mutex
	group := errgroup.Group{}

	for _, entry := range p.Tracks {
		entry := entry
		group.Go(func() error {
			return s.Add(p.Dir, entry, &mu)
		})
	}

	err = group.Wait()
	return s, err
}
