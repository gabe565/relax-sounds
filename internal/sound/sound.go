package sound

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
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

var ErrInvalidMetaFileType = errors.New("invalid meta file type")

func Load(fsys fs.FS, path string) (sound Sound, err error) {
	path = filepath.Join("/", path)[1:]

	f, err := fsys.Open(path)
	if err != nil {
		return sound, err
	}
	defer func(f fs.File) {
		_ = f.Close()
	}(f)

	ext := filepath.Ext(path)

	switch ext {
	case ".yaml", ".yml":
		err = yaml.NewDecoder(f).Decode(&sound)
	case ".json":
		err = json.NewDecoder(f).Decode(&sound)
	default:
		return sound, fmt.Errorf("%w: %#v", ErrInvalidMetaFileType, path)
	}
	if err != nil {
		return sound, err
	}

	rawPath := strings.TrimPrefix(path, "meta")[1:]
	if sound.Filename == "" {
		sound.Filename = strings.TrimSuffix(rawPath, ext) + ".ogg"
	}

	if sound.Id == "" {
		id := strings.TrimSuffix(filepath.Base(path), ext)
		sound.Id = id
	}

	return sound, err
}
