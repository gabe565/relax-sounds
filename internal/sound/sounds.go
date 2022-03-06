package sound

import (
	"io/fs"
)

func LoadAll(fsys fs.FS) (sounds []Sound, err error) {
	err = fs.WalkDir(fsys, "meta", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		sound, err := Load(fsys, path)
		if err != nil {
			return err
		}

		sounds = append(sounds, sound)
		return nil
	})
	return sounds, err
}
