package preset

type Track struct {
	Id     string
	Volume *float64
	Rate   *float64
}

func (t Track) GetVolume() float64 {
	if t.Volume == nil {
		return 1
	}
	return *t.Volume
}

func (t Track) GetRate() float64 {
	if t.Rate == nil {
		return 1
	}
	rate := *t.Rate
	if rate < 0.5 {
		return 0.5
	}
	if rate > 1.5 {
		return 1.5
	}
	return rate
}

type Preset struct {
	Tracks []Track
}

func (preset *Preset) Add(tracks ...Track) {
	preset.Tracks = append(preset.Tracks, tracks...)
}
