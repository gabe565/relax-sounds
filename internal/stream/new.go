package stream

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"gabe565.com/relax-sounds/internal/config"
	"gabe565.com/relax-sounds/internal/preset"
	"github.com/pocketbase/pocketbase/core"
	"golang.org/x/sync/errgroup"
)

var ErrInvalidRecordID = errors.New("invalid record ID")

func New(conf *config.Config, p preset.Preset) (Streams, error) {
	s := make(Streams, 0, len(p))

	var mu sync.Mutex
	group := errgroup.Group{}
	group.SetLimit(5)

	ids := make([]string, 0, len(p))
	for _, entry := range p {
		ids = append(ids, entry.ID)
	}
	records, err := conf.App.FindRecordsByIds("sounds", ids)
	if err != nil {
		return s, err
	}

	storageDir := filepath.Join(conf.App.DataDir(), "storage")
	for _, entry := range p {
		var record *core.Record
		for _, v := range records {
			if v.Id == entry.ID {
				record = v
				break
			}
		}
		if record == nil {
			return s, fmt.Errorf("%w: %s", ErrInvalidRecordID, entry.ID)
		}

		group.Go(func() error {
			files := record.GetStringSlice("file")
			var preferredFile string
			for _, preferredExt := range preferredExts() {
				i := slices.IndexFunc(files, func(s string) bool {
					return filepath.Ext(s) == preferredExt
				})
				if i != -1 {
					preferredFile = files[i]
					break
				}
			}
			if preferredFile == "" {
				preferredFile = files[0]
			}
			path := filepath.Join(storageDir, record.BaseFilesPath(), preferredFile)

			f, err := os.Open(path)
			if err != nil {
				return err
			}

			return s.Add(conf, f, entry, &mu)
		})
	}

	err = group.Wait()
	return s, err
}

func preferredExts() []string {
	return []string{".wav", ".ogg"}
}
