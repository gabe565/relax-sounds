package stream

import (
	"errors"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
	"github.com/gabe565/relax-sounds/internal/util"
	"io"
	"io/fs"
)

type File struct {
	fs.File
}

var ErrUnsupportedFileType = errors.New("unsupported file type")

func (f File) Decode() (beep.StreamSeekCloser, beep.Format, error) {
	contentType, err := util.GetTypeFromFile(f)
	if err != nil {
		return nil, beep.Format{}, err
	}

	_, err = f.File.(io.ReadSeeker).Seek(0, io.SeekStart)
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
