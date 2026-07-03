package stream

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"gabe565.com/relax-sounds/internal/config"
	"gabe565.com/relax-sounds/internal/handlers/mix/preset"
	"github.com/pocketbase/pocketbase/core"
)

var ErrInvalidRecordID = errors.New("invalid record ID")

func New(conf *config.Config, p preset.Preset) (Streams, error) {
	s := make(Streams, 0, len(p))

	ids := make([]string, 0, len(p))
	for _, entry := range p {
		ids = append(ids, entry.ID)
	}
	records, err := conf.App.FindRecordsByIds("sounds", ids)
	if err != nil {
		return s, err
	}

	recordByID := make(map[string]*core.Record, len(records))
	for _, r := range records {
		recordByID[r.Id] = r
	}

	storageDir := filepath.Join(conf.App.DataDir(), "storage")
	for _, entry := range p {
		record, ok := recordByID[entry.ID]
		if !ok {
			return s, fmt.Errorf("%w: %s", ErrInvalidRecordID, entry.ID)
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
