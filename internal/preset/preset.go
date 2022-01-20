package preset

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

type Track struct {
	Key    string
	Volume float64
}

func (track Track) ToShorthand() ShorthandTrack {
	return ShorthandTrack{track.Key, track.Volume}
}

func (track Track) Path() string {
	return filepath.Clean(fmt.Sprintf("audio/%s.ogg", track.Key))
}

type Preset struct {
	Dir    fs.FS
	Tracks []Track
}

func (preset Preset) ToShorthand() PresetShorthand {
	shorthand := PresetShorthand{}
	for _, track := range preset.Tracks {
		shorthand = append(shorthand, track.ToShorthand())
	}
	return shorthand
}

func (preset *Preset) Add(tracks ...Track) {
	preset.Tracks = append(preset.Tracks, tracks...)
}
