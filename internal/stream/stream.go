package stream

import (
	"github.com/gabe565/relax-sounds/internal/preset"
	"github.com/pocketbase/pocketbase/daos"
	"golang.org/x/sync/errgroup"
	"io/fs"
	"path/filepath"
	"sync"
)

func New(dataDir fs.FS, dao *daos.Dao, p preset.Preset) (stream Streams, err error) {
	s := make(Streams, 0, len(p.Tracks))

	var mu sync.Mutex
	group := errgroup.Group{}

	for _, entry := range p.Tracks {
		entry := entry
		group.Go(func() error {
			record, err := dao.FindRecordById("sounds", entry.Id)
			if err != nil {
				return err
			}

			path := filepath.Join(record.BaseFilesPath(), record.Get("file").(string))

			f, err := dataDir.Open(path)
			if err != nil {
				return err
			}

			return s.Add(f, entry, &mu)
		})
	}

	err = group.Wait()
	return s, err
}
