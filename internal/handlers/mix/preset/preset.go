package preset

import "strconv"

//nolint:recvcheck
type Track struct {
	ID      string `json:"id"`
	ShortID *int   `json:"short_id"`

	Volume *float64 `json:"volume"`
	Rate   *float64 `json:"rate"`
	Pan    *float64 `json:"pan"`
}

func (t Track) Ref() string {
	if t.ShortID != nil {
		return "#" + strconv.Itoa(*t.ShortID)
	}
	return t.ID
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
	return max(0.5, min(*t.Rate, 1.5))
}

func (t Track) GetPan() float64 {
	if t.Pan == nil {
		return 0
	}
	return max(-1, min(*t.Pan, 1))
}

type Preset []Track
