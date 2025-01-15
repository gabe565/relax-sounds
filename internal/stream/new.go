package stream

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"gabe565.com/relax-sounds/internal/config"
	"gabe565.com/relax-sounds/internal/preset"
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
