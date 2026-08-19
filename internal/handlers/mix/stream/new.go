package stream

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"gabe565.com/relax-sounds/internal/config"
	"gabe565.com/relax-sounds/internal/handlers/mix/preset"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

var ErrInvalidRecordID = errors.New("invalid record ID")

func New(conf *config.Config, p preset.Preset) (Streams, error) {
	s := make(Streams, 0, len(p))

	recordByID, recordByShortID, err := findRecords(conf, p)
	if err != nil {
		return s, err
	}

	storageDir := filepath.Join(conf.App.DataDir(), "storage")
	for _, entry := range p {
		var record *core.Record
		var ok bool
		if entry.ShortID != nil {
			record, ok = recordByShortID[*entry.ShortID]
		} else {
			record, ok = recordByID[entry.ID]
		}
		if !ok {
			return s, fmt.Errorf("%w: %s", ErrInvalidRecordID, entry.Ref())
		}

		files := record.GetStringSlice("file")
		path := filepath.Join(storageDir, record.BaseFilesPath(), preferredFile(files))

		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		if err := s.Add(conf, f, entry); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func findRecords(conf *config.Config, p preset.Preset) (map[string]*core.Record, map[int]*core.Record, error) {
	ids := make([]string, 0, len(p))
	shortIDs := make([]any, 0, len(p))

	for _, entry := range p {
		if entry.ShortID != nil {
			shortIDs = append(shortIDs, *entry.ShortID)
		} else {
			ids = append(ids, entry.ID)
		}
	}

	recordByID := make(map[string]*core.Record, len(ids))
	if len(ids) != 0 {
		records, err := conf.App.FindRecordsByIds("sounds", ids)
		if err != nil {
			return nil, nil, err
		}
		for _, r := range records {
			recordByID[r.Id] = r
		}
	}

	recordByShortID := make(map[int]*core.Record, len(shortIDs))
	if len(shortIDs) != 0 {
		var records []*core.Record
		if err := conf.App.RecordQuery("sounds").
			AndWhere(dbx.In("short_id", shortIDs...)).
			All(&records); err != nil {
			return nil, nil, err
		}
		for _, r := range records {
			recordByShortID[r.GetInt("short_id")] = r
		}
	}

	return recordByID, recordByShortID, nil
}

func preferredExts() []string {
	return []string{".wav", ".ogg"}
}

func preferredFile(files []string) string {
	for _, ext := range preferredExts() {
		i := slices.IndexFunc(files, func(s string) bool {
			return filepath.Ext(s) == ext
		})
		if i != -1 {
			return files[i]
		}
	}
	return files[0]
}
