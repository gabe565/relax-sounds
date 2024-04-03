package stream

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"sync"

	"github.com/gabe565/relax-sounds/internal/preset"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"golang.org/x/sync/errgroup"
)

var ErrInvalidRecordID = errors.New("invalid record ID")

func New(dataDir fs.FS, dao *daos.Dao, p preset.Preset) (Streams, error) {
	s := make(Streams, 0, len(p.Tracks))

	var mu sync.Mutex
	group := errgroup.Group{}
	group.SetLimit(5)

	ids := make([]string, 0, len(p.Tracks))
	for _, entry := range p.Tracks {
		ids = append(ids, entry.ID)
	}
	records, err := dao.FindRecordsByIds("sounds", ids)
	if err != nil {
		return s, err
	}

	for _, entry := range p.Tracks {
		var record *models.Record
		for _, v := range records {
			if v.GetId() == entry.ID {
				record = v
				break
			}
		}
		if record == nil {
			return s, fmt.Errorf("%w: %s", ErrInvalidRecordID, entry.ID)
		}

		group.Go(func() error {
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
