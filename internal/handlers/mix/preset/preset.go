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
	switch {
	case t.Rate == nil:
		return 1
	case *t.Rate < 0.5:
		return 0.5
	case *t.Rate > 1.5:
		return 1.5
	default:
		return *t.Rate
	}
}

func (t Track) GetPan() float64 {
	switch {
	case t.Pan == nil:
		return 0
	case *t.Pan < -1:
		return -1
	case *t.Pan > 1:
		return 1
	default:
		return *t.Pan
	}
}

type Preset []Track
