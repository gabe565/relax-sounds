package playlist

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

type Track struct {
	Key    string
	Volume float64
}

func (entry Track) ToShorthand() ShorthandTrack {
	return ShorthandTrack{entry.Key, entry.Volume}
}

func (entry Track) Path() string {
	return filepath.Clean(fmt.Sprintf("audio/%s.ogg", entry.Key))
}

type Playlist struct {
	Dir    fs.FS
	Tracks []Track
}

func (playlist Playlist) ToShorthand() PlaylistShorthand {
	shorthand := PlaylistShorthand{}
	for _, track := range playlist.Tracks {
		shorthand = append(shorthand, track.ToShorthand())
	}
	return shorthand
}

func (playlist *Playlist) Add(tracks ...Track) {
	playlist.Tracks = append(playlist.Tracks, tracks...)
}
