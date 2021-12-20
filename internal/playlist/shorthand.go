package playlist

import (
	"errors"
	"fmt"
)

var ErrInvalidShorthand = errors.New("invalid shorthand")

type ShorthandEntry []interface{}

type Shorthand []ShorthandEntry

func (shorthand Shorthand) ToPlaylist() (Playlist, error) {
	playlist := make(Playlist, 0, len(shorthand))
	for _, value := range shorthand {
		if len(value) != 2 {
			// Shorthand is not a tuple
			return playlist, ErrInvalidShorthand
		}
		playlist = append(playlist, Entry{
			Key:    fmt.Sprintf("%v", value[0]),
			Volume: value[1].(float64) - 1,
		})
	}
	return playlist, nil
}
