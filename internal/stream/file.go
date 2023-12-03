package stream

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/gabe565/relax-sounds/internal/util"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/flac"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/vorbis"
	"github.com/gopxl/beep/wav"
)

var ErrUnsupportedFileType = errors.New("unsupported file type")

func Decode(f fs.File) (beep.StreamSeekCloser, beep.Format, error) {
	contentType, err := util.GetTypeFromFile(f)
	if err != nil {
		return nil, beep.Format{}, err
	}

	switch contentType {
	case "application/ogg", "audio/ogg":
		return vorbis.Decode(f)
	case "audio/mpeg":
		return mp3.Decode(f)
	case "audio/wave", "audio/x-wav":
		return wav.Decode(f)
	case "audio/x-flac":
		return flac.Decode(f)
	}
	return nil, beep.Format{}, fmt.Errorf("%w: %s", ErrUnsupportedFileType, contentType)
}
