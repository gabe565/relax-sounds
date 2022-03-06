package sound

import (
	"errors"
	"golang.org/x/sync/errgroup"
	"io/fs"
	"log"
	"sync"
)

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
	return sounds, err
}
