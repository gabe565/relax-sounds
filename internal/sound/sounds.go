package sound

import (
	"errors"
	"io/fs"
	"log"
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
			if errors.Is(err, ErrInvalidMetaFileType) {
				log.Println("WARN: " + err.Error())
				return nil
			} else {
				return err
			}
		}

		sounds = append(sounds, sound)
		return nil
	})
	return sounds, err
}
