package preset

import "path/filepath"

type Track struct {
	Id     string
	Volume float64
}

func (track Track) ToShorthand() ShorthandTrack {
	return ShorthandTrack{track.Id, track.Volume}
}

func (track Track) Path() string {
	return filepath.Join("/", "audio", track.Id+".ogg")[1:]
}
