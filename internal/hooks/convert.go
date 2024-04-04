package hooks

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/rs/zerolog/log"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func Convert(app *pocketbase.PocketBase) func(e *core.ModelEvent) error {
	dataDir := filepath.Join(app.DataDir(), "storage")

	var skipConvert bool
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		skipConvert = true
		log.Warn().Err(err).Msg("ffmpeg is required for secondary audio file creation")
	}

	return func(e *core.ModelEvent) error {
		record := e.Model.(*models.Record)

		files := record.GetStringSlice("file")
		if len(files) == 0 || len(files) > 1 {
			return nil
		}

		ext := filepath.Ext(files[0])
		if ext != ".ogg" {
			return nil
		}

		if skipConvert {
			log.Error().Msg("ffmpeg is required for secondary audio file creation")
			return nil
		}

		path := filepath.Join(dataDir, record.BaseFilesPath(), files[0])
		dstPath := strings.TrimSuffix(files[0], ext) + ".mp3"
		log.Info().
			Str("src", filepath.Base(path)).
			Str("dst", filepath.Base(dstPath)).
			Msg("creating secondary audio file")

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
			log.Err(err).Msg("failed to convert file")
			return err
		}

		form := forms.NewRecordUpsert(app, record)
		file, err := filesystem.NewFileFromBytes(dst.Bytes(), dstPath)
		file.Name = dstPath
		if err != nil {
			return err
		}
		if err := form.AddFiles("file", file); err != nil {
			return err
		}

		return form.Submit()
	}
}

func ConvertAll(app *pocketbase.PocketBase) error {
	start := time.Now()

	records, err := app.Dao().FindRecordsByExpr("sounds")
	if err != nil {
		return err
	}

	for _, record := range records {
		form := forms.NewRecordUpsert(app, record)
		if err := form.Submit(); err != nil {
			return err
		}
	}

	log.Info().Str("took", time.Since(start).Truncate(time.Second).String()).
		Msg("finished creating secondary audio files")
	return nil
}
