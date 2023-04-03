package sound

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"io/fs"
	"regexp"
	"sort"
	"strconv"
)

var sortMixed = regexp.MustCompile(`^\d+`)

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
				log.Warn(err.Error())
				return nil
			} else {
				return err
			}
		}

		sounds = append(sounds, sound)
		return nil
	})
	if err != nil {
		return sounds, err
	}

	// Sort numerically if possible. Falls back to string sort
	sort.Slice(sounds, func(i, j int) bool {
		left := sounds[i].Id
		right := sounds[j].Id

		if leftInt, err := strconv.Atoi(sortMixed.FindString(left)); err == nil {
			if rightInt, err := strconv.Atoi(sortMixed.FindString(right)); err == nil {
				if leftInt != rightInt {
					// Compare ints
					return leftInt < rightInt
				}
			}
		}

		// Compare strings
		return left < right
	})

	return sounds, err
}
