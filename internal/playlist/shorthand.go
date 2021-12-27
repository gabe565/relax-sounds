package playlist

import (
	"errors"
	"fmt"
	"path/filepath"
)

var ErrInvalidShorthand = errors.New("invalid shorthand")

type ShorthandTrack [2]interface{}

type PlaylistShorthand []ShorthandTrack

func (shorthand PlaylistShorthand) ToPlaylist(staticDir string) (Playlist, error) {
	playlist := Playlist{}
	for _, value := range shorthand {
		track := Track{
			Key:    fmt.Sprintf("%v", value[0]),
			Path:   fmt.Sprintf("%v/public/audio/%v.ogg", staticDir, value[0]),
			Volume: value[1].(float64) - 1,
		}
		track.Path = filepath.Join(staticDir+"/audio", filepath.Clean("/"+track.Key+".ogg"))
		playlist.Add(track)
	}
	return playlist, nil
}
