package preset

import (
	"errors"
	"fmt"
	"io/fs"
)

var ErrInvalidShorthand = errors.New("invalid shorthand")

type ShorthandTrack [2]interface{}

type PresetShorthand []ShorthandTrack

func (shorthand PresetShorthand) ToPreset(dataDir fs.FS) (Preset, error) {
	preset := Preset{
		Dir: dataDir,
	}
	for _, value := range shorthand {
		track := Track{
			Key:    fmt.Sprintf("%v", value[0]),
			Volume: value[1].(float64),
		}
		preset.Add(track)
	}
	return preset, nil
}
