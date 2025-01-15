package hooks

import (
	"errors"
	"fmt"
	"io"
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
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func altExts() []string {
	return []string{".mp3", ".wav"}
}

func Convert(app *pocketbase.PocketBase) func(e *core.ModelEvent) error {
	dataDir := filepath.Join(app.DataDir(), "storage")

	var skipConvert bool
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		skipConvert = true
		slog.Warn("ffmpeg is required for secondary audio file creation", "error", err)
	}

	return func(e *core.ModelEvent) error {
		record := e.Model.(*core.Record)

		files := record.GetStringSlice("file")
		exts := make([]string, 0, len(files))
		for _, file := range files {
			exts = append(exts, filepath.Ext(file))
		}

		if skipConvert {
			slog.Error("ffmpeg is required for alt audio file creation")
			return nil
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

			var kwargs ffmpeg.KwArgs
			if altExt == ".mp3" {
				kwargs = ffmpeg.KwArgs{"qscale:a": "2"}
			}

			var errBuf strings.Builder
			err = ffmpeg.Input(path).
				Output(dstPath, kwargs).
				WithOutput(io.Discard, &errBuf).
				Silent(true).
				Run()
			if err != nil {
				err := fmt.Errorf("%w: %s", err, errBuf.String())
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
