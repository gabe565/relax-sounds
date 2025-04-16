package encoder

import (
	"errors"
	"fmt"
	"io"

	"gabe565.com/relax-sounds/internal/config"
	"github.com/gopxl/beep/v2"
)

//go:generate go tool enumer -type Format -trimprefix Format -transform lower -output format_string.go

type Format uint8

const (
	FormatMP3 Format = iota
)

var ErrInvalidFormat = errors.New("invalid file format")

func (i Format) ContentType() string {
	if i == FormatMP3 {
		return "audio/mp3"
	}
	return ""
}

func (i Format) NewEncoder(conf *config.Config, w io.Writer, format beep.Format) (Encoder, error) {
	if i == FormatMP3 {
		return NewEncoder(conf, w, format)
	}
	return nil, fmt.Errorf("%w: %s", ErrInvalidFormat, i)
}
