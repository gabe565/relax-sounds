package playlist

import (
	"errors"
	"fmt"
	"io/fs"
)

var ErrInvalidShorthand = errors.New("invalid shorthand")

type ShorthandTrack [2]interface{}

type PlaylistShorthand []ShorthandTrack

func (shorthand PlaylistShorthand) ToPlaylist(dataDir fs.FS) (Playlist, error) {
	playlist := Playlist{
		Dir: dataDir,
	}
	for _, value := range shorthand {
		track := Track{
			Key:    fmt.Sprintf("%v", value[0]),
			Volume: value[1].(float64) - 1,
		}
		playlist.Add(track)
	}
	return playlist, nil
}
