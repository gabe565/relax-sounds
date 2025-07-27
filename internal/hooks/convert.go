package hooks

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

func altExts() []string {
	return []string{".mp3", ".wav"}
}

func Convert(app *pocketbase.PocketBase) func(e *core.ModelEvent) error {
	dataDir := filepath.Join(app.DataDir(), "storage")

	if _, err := exec.LookPath("ffmpeg"); err != nil {
		slog.Warn("ffmpeg is required for alt audio file creation", "error", err)
	}

	return func(e *core.ModelEvent) error {
		if _, err := exec.LookPath("ffmpeg"); err != nil {
			slog.Warn("ffmpeg is required for alt audio file creation", "error", err)
			return nil
		}

		record := e.Model.(*core.Record) //nolint:errcheck

		files := record.GetStringSlice("file")
		exts := make([]string, 0, len(files))
		for _, file := range files {
			exts = append(exts, filepath.Ext(file))
		}

		tmpDir, err := os.MkdirTemp("", "relax-sounds-*")
		if err != nil {
			return err
		}
		defer func() { _ = os.RemoveAll(tmpDir) }()

		var changed bool
		for _, altExt := range altExts() {
			if slices.Contains(exts, altExt) {
				continue
			}

			path := filepath.Join(dataDir, record.BaseFilesPath(), files[0])
			dstPath := filepath.Join(tmpDir, strings.TrimSuffix(files[0], exts[0])+altExt)
			slog.Info("Creating alt audio file",
				"src", filepath.Base(path),
				"dst", filepath.Base(dstPath),
			)

			args := []string{
				"-hide_banner",
				"-loglevel", "error",
				"-i", path,
			}
			if altExt == ".mp3" {
				args = append(args, "-qscale:a", "2")
			}
			args = append(args, dstPath)

			cmd := exec.CommandContext(e.Context, "ffmpeg", args...)
			b, err := cmd.CombinedOutput()
			if err != nil {
				err := fmt.Errorf("%w: %s", err, b)
				slog.Error("Failed to convert audio file", "error", err)
				return err
			}

			file, err := filesystem.NewFileFromPath(dstPath)
			if err != nil {
				return err
			}
			file.Name = filepath.Base(dstPath)
			record.Set("file+", file)
			changed = true
		}

		if changed {
			if err := app.Save(record); err != nil {
				return err
			}
		}

		return e.Next()
	}
}

func ConvertAll(app core.App) error {
	records, err := app.FindAllRecords("sounds")
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return nil
	}

	var errs []error
	for _, record := range records {
		files := record.GetStringSlice("file")
		for _, altExt := range altExts() {
			i := slices.IndexFunc(files, func(s string) bool {
				return filepath.Ext(s) == altExt
			})
			if i == -1 {
				form := forms.NewRecordUpsert(app, record)
				if err := form.Submit(); err != nil {
					errs = append(errs, err)
				}
				break
			}
		}
	}
	return errors.Join(errs...)
}
