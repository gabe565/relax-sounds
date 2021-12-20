package stream

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/vorbis"
	"github.com/gabe565/relax-sounds/internal/playlist"
	"os"
	"path/filepath"
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

func (stream *Stream) Add(entry playlist.Entry) error {
	audiopath := filepath.Join("public/audio", filepath.Clean("/"+entry.Key+".ogg"))
	infile, err := os.Open(audiopath)
	if err != nil {
		return err
	}

	streamer, format, err := vorbis.Decode(infile)
	if err != nil {
		return err
	}

	volumeStreamer := &effects.Volume{
		Streamer: beep.Loop(-1, streamer),
		Base:     4,
		Volume:   entry.Volume,
		Silent:   false,
	}
	stream.closers = append(stream.closers, streamer)
	stream.streamers = append(stream.streamers, volumeStreamer)
	stream.Formats = append(stream.Formats, format)
	return nil
}

func (stream *Stream) Mix() beep.Streamer {
	return beep.Mix(stream.streamers...)
}

func New(plist playlist.Playlist) (stream Stream, err error) {
	s := Stream{
		closers:   make([]beep.StreamSeekCloser, 0, len(plist)),
		streamers: make([]beep.Streamer, 0, len(plist)),
		Formats:   make([]beep.Format, 0, len(plist)),
	}
	for _, entry := range plist {
		if err = s.Add(entry); err != nil {
			return s, err
		}
	}
	return s, nil
}
