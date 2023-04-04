package preset

type Preset struct {
	Tracks []Track
}

func (preset Preset) ToShorthand() Shorthand {
	shorthand := Shorthand{}
	for _, track := range preset.Tracks {
		shorthand = append(shorthand, track.ToShorthand())
	}
	return shorthand
}

func (preset *Preset) Add(tracks ...Track) {
	preset.Tracks = append(preset.Tracks, tracks...)
}
