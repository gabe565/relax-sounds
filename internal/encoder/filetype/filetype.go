package filetype

import (
	"errors"
	"fmt"
	"io"

	"gabe565.com/relax-sounds/internal/encoder"
	"gabe565.com/relax-sounds/internal/encoder/mp3"
	"github.com/gopxl/beep/v2"
)

//go:generate go run github.com/dmarkham/enumer -type FileType -transform lower -output filetype_string.go

type FileType uint8

const (
	MP3 FileType = iota
)

var ErrInvalidFileType = errors.New("invalid file type")

func (i FileType) ContentType() string {
	if i == MP3 {
		return "audio/mp3"
	}
	return ""
}

func (i FileType) NewEncoder(w io.Writer, format beep.Format) (encoder.Encoder, error) {
	if i == MP3 {
		return mp3.NewEncoder(w, format)
	}
	return nil, fmt.Errorf("%w: %s", ErrInvalidFileType, i)
}
