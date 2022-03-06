package preset

import "path/filepath"

type Track struct {
	Key    string
	Volume float64
}

func (track Track) ToShorthand() ShorthandTrack {
	return ShorthandTrack{track.Key, track.Volume}
}

func (track Track) Path() string {
	return filepath.Join("/", "audio", track.Key+".ogg")[1:]
}
