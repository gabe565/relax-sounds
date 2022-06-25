package handlers

import (
	"encoding/json"
	"github.com/gabe565/relax-sounds/internal/sound"
	"io/fs"
	"net/http"
)

func Sounds(fsys fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sounds, err := sound.LoadAll(fsys)
		if err != nil {
			panic(err)
		}

		soundJson, err := json.Marshal(sounds)
		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(soundJson)
		if err != nil {
			panic(err)
		}
	}
}
