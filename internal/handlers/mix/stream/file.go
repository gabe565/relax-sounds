package stream

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/flac"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/vorbis"
	"github.com/gopxl/beep/v2/wav"
)

var ErrUnsupportedFileType = errors.New("unsupported file type")

func Decode(f fs.File) (beep.StreamSeekCloser, beep.Format, error) {
	stat, err := f.Stat()
	if err != nil {
		return nil, beep.Format{}, err
	}

	ext := filepath.Ext(stat.Name())

	switch ext {
	case ".ogg":
		return vorbis.Decode(f)
	case ".mp3":
		return mp3.Decode(f)
	case ".wav":
		return wav.Decode(f)
	case ".flac":
		return flac.Decode(f)
	}
	return nil, beep.Format{}, fmt.Errorf("%w: %s", ErrUnsupportedFileType, ext)
}
