package sound

import (
	"errors"
	"golang.org/x/sync/errgroup"
	"io/fs"
	"log"
	"regexp"
	"sort"
	"strconv"
	"sync"
)

var sortMixed = regexp.MustCompile("^[0-9]+")

func LoadAll(fsys fs.FS) (sounds []Sound, err error) {
	var mu sync.Mutex
	group := errgroup.Group{}

	err = fs.WalkDir(fsys, "meta", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		group.Go(func() error {
			sound, err := Load(fsys, path)
			if err != nil {
				if errors.Is(err, ErrInvalidMetaFileType) {
					log.Println("WARN: " + err.Error())
					return nil
				} else {
					return err
				}
			}
			mu.Lock()
			defer mu.Unlock()
			sounds = append(sounds, sound)
			return nil
		})
		return nil
	})

	err = group.Wait()

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
