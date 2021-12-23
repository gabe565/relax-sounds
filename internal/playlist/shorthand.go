package playlist

import (
	"errors"
	"fmt"
)

var ErrInvalidShorthand = errors.New("invalid shorthand")

type ShorthandTrack [2]interface{}

type PlaylistShorthand []ShorthandTrack

func (shorthand PlaylistShorthand) ToPlaylist() (Playlist, error) {
	playlist := Playlist{}
	for _, value := range shorthand {
		playlist.Add(Track{
			Key:    fmt.Sprintf("%v", value[0]),
			Volume: value[1].(float64) - 1,
		})
	}
	return playlist, nil
}
