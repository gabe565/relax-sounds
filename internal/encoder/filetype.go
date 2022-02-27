package encoder

import (
	"errors"
	"fmt"
	"github.com/faiface/beep"
	"github.com/gabe565/relax-sounds/internal/encoder/mp3"
	"github.com/gabe565/relax-sounds/internal/encoder/wav"
	"io"
)

//go:generate stringer -type FileType -linecomment

type FileType uint8

const (
	FileTypeWav FileType = iota // wav
	FileTypeMp3                 // mp3
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
	case FileTypeWav:
		return "audio/wav"
	case FileTypeMp3:
		return "audio/mp3"
	}
	return ""
}

func (i FileType) NewEncoder(w io.Writer, format beep.Format) (Encoder, error) {
	switch i {
	case FileTypeWav:
		return wav.NewEncoder(w, format)
	case FileTypeMp3:
		return mp3.NewEncoder(w, format)
	}
	return nil, fmt.Errorf("%w: %s", ErrInvalidFileType, i)
}
