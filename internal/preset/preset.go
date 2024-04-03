package preset

type Track struct {
	ID     string   `json:"id"`
	Volume *float64 `json:"volume"`
	Rate   *float64 `json:"rate"`
	Pan    *float64 `json:"pan"`
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

func (t Track) GetPan() float64 {
	if t.Pan == nil {
		return 0
	}
	pan := *t.Pan
	if pan < -1 {
		return -1
	}
	if pan > 1 {
		return 1
	}
	return pan
}

type Preset struct {
	Tracks []Track `json:"tracks"`
}

func (preset *Preset) Add(tracks ...Track) {
	preset.Tracks = append(preset.Tracks, tracks...)
}
