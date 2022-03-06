package preset

import (
	"errors"
	"fmt"
	"io/fs"
)

var ErrInvalidId = errors.New("invalid id")
var ErrInvalidVolume = errors.New("invalid volume")

type ShorthandTrack [2]interface{}

type PresetShorthand []ShorthandTrack

func (shorthand PresetShorthand) ToPreset(dataDir fs.FS) (Preset, error) {
	preset := Preset{
		Dir: dataDir,
	}
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
