package stream

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/vorbis"
	"github.com/gabe565/relax-sounds/internal/playlist"
	"golang.org/x/sync/errgroup"
	"os"
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

func (stream *Stream) Add(entry playlist.Track, mu *sync.Mutex) error {
	infile, err := os.Open(entry.Path)
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

func New(plist playlist.Playlist) (stream Stream, err error) {
	s := Stream{
		closers:   make([]beep.StreamSeekCloser, 0, len(plist.Tracks)),
		streamers: make([]beep.Streamer, 0, len(plist.Tracks)),
		Formats:   make([]beep.Format, 0, len(plist.Tracks)),
	}

	var mu sync.Mutex
	group := errgroup.Group{}

	for _, entry := range plist.Tracks {
		entry := entry
		group.Go(func() error {
			return s.Add(entry, &mu)
		})
	}

	err = group.Wait()
	return s, err
}
