package preset

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidId     = errors.New("invalid id")
	ErrInvalidVolume = errors.New("invalid volume")
)

type ShorthandTrack [2]any

type Shorthand []ShorthandTrack

func (shorthand Shorthand) ToPreset() (Preset, error) {
	var preset Preset
	for _, value := range shorthand {
		id, ok := value[0].(string)
		if !ok {
			return preset, fmt.Errorf("%w: %#v", ErrInvalidId, value[0])
		}

		volume, ok := value[1].(float64)
		if !ok {
			return preset, fmt.Errorf("%w: %#v", ErrInvalidVolume, value[1])
		}

		track := Track{
			Id:     id,
			Volume: volume,
		}
		preset.Add(track)
	}
	return preset, nil
}
