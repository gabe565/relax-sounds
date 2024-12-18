package hooks

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

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
		if len(files) != 1 {
			return nil
		}

		ext := filepath.Ext(files[0])
		if ext != ".ogg" {
			return nil
		}

		if skipConvert {
			slog.Error("ffmpeg is required for secondary audio file creation")
			return nil
		}

		path := filepath.Join(dataDir, record.BaseFilesPath(), files[0])
		dstPath := strings.TrimSuffix(files[0], ext) + ".mp3"
		slog.Info("Creating secondary audio file",
			"src", filepath.Base(path),
			"dst", filepath.Base(dstPath),
		)

		var dst bytes.Buffer
		var errBuf strings.Builder
		err := ffmpeg.Input(path).
			Output("pipe:", ffmpeg.KwArgs{
				"format":   "mp3",
				"codec:a":  "libmp3lame",
				"qscale:a": "2",
			}).
			WithOutput(&dst, &errBuf).
			Silent(true).
			Run()
		if err != nil {
			err := fmt.Errorf("%w: %s", err, errBuf.String())
			slog.Error("Failed to convert audio file", "error", err)
			return err
		}

		file, err := filesystem.NewFileFromBytes(dst.Bytes(), dstPath)
		if err != nil {
			return err
		}
		file.Name = dstPath
		record.Set("file+", file)

		if err := app.Save(record); err != nil {
			return err
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
		if len(record.GetStringSlice("file")) == 1 {
			form := forms.NewRecordUpsert(app, record)
			if err := form.Submit(); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errors.Join(errs...)
}
