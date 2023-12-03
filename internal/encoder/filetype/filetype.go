package filetype

import (
	"errors"
	"fmt"
	"io"

	"github.com/gabe565/relax-sounds/internal/encoder"
	"github.com/gabe565/relax-sounds/internal/encoder/mp3"
	"github.com/gopxl/beep"
)

//go:generate stringer -type FileType -linecomment

type FileType uint8

const (
	Mp3 FileType = iota // mp3
)

var ErrInvalidFileType = errors.New("invalid file type")

func (i FileType) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

func (i *FileType) UnmarshalText(b []byte) error {
	s := string(b)
	for j := 0; j < len(_FileType_index)-1; j++ {
		if s == _FileType_name[_FileType_index[j]:_FileType_index[j+1]] {
			*i = FileType(j)
			return nil
		}
	}
	return fmt.Errorf("%w: %s", ErrInvalidFileType, s)
}

func (i FileType) ContentType() string {
	switch i {
	case Mp3:
		return "audio/mp3"
	}
	return ""
}

func (i FileType) NewEncoder(w io.Writer, format beep.Format) (encoder.Encoder, error) {
	switch i {
	case Mp3:
		return mp3.NewEncoder(w, format)
	}
	return nil, fmt.Errorf("%w: %s", ErrInvalidFileType, i)
}
