package tag

import (
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"path/filepath"
)

type Tag struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon,omitempty"`
}

var ErrNoTagFile = errors.New("tag config not found")

var tagFiles = []string{"tags.yaml", "tags.yml", "tags.json"}

func LoadAll(fsys fs.FS) (tags []Tag, err error) {
	var f fs.File
	var path string
	for _, path = range tagFiles {
		f, err = fsys.Open(path)
		if err == nil {
			// Found file
			break
		} else {
			// Continue on file not found
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return tags, err
		}
	}
	if f == nil {
		return tags, ErrNoTagFile
	}
	defer func(f fs.File) {
		_ = f.Close()
	}(f)

	ext := filepath.Ext(path)[1:]
	switch ext {
	case "yaml", "yml":
		err = yaml.NewDecoder(f).Decode(&tags)
	case "json":
		err = json.NewDecoder(f).Decode(&tags)
	}

	return tags, err
}
