package mp3

import (
	"io"

	"gabe565.com/relax-sounds/internal/config"
	"github.com/gopxl/beep/v2"
	"github.com/viert/go-lame"
)

type Encoder struct {
	Encoder *lame.Encoder
	Format  beep.Format
}

func (e Encoder) Write(p []byte) (int, error) {
	return e.Encoder.Write(p)
}

func (e Encoder) Close() error {
	e.Encoder.Close()
	return nil
}

func NewEncoder(conf *config.Config, w io.Writer, format beep.Format) (io.WriteCloser, error) {
	var err error
	enc := Encoder{
		Encoder: lame.NewEncoder(w),
		Format:  format,
	}
	if err = enc.Encoder.SetVBR(lame.VBRDefault); err != nil {
		return enc, err
	}
	if err = enc.Encoder.SetVBRQuality(conf.LAMEQuality); err != nil {
		return enc, err
	}
	return enc, nil
}
