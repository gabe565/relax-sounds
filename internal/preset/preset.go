package preset

import (
	"io/fs"
)

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
