package sound

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

type Sound struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Icon     string   `json:"icon"`
	Filename string   `json:"filename"`
	Tags     []string `json:"tags"`
}

func Load(fsys fs.FS, path string) (sound Sound, err error) {
	f, err := fsys.Open(path)
	if err != nil {
		return sound, err
	}
	defer func(f fs.File) {
		_ = f.Close()
	}(f)

	err = json.NewDecoder(f).Decode(&sound)
	if err != nil {
		return sound, err
	}

	ext := filepath.Ext(path)
	rawPath := strings.TrimPrefix(path, fmt.Sprintf("meta%c", filepath.Separator))
	if sound.Filename == "" {
		sound.Filename = strings.TrimSuffix(rawPath, ext) + ".ogg"
	}

	if sound.Id == "" {
		id := strings.TrimSuffix(filepath.Base(path), ext)
		sound.Id = id
	}

	return sound, err
}
