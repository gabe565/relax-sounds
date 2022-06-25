package handlers

import (
	"encoding/json"
	"github.com/gabe565/relax-sounds/internal/tag"
	"io/fs"
	"net/http"
)

func Tags(fsys fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tags, err := tag.LoadAll(fsys)
		if err != nil {
			panic(err)
		}

		soundJson, err := json.Marshal(tags)
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
