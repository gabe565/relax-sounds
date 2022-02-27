package mp3

import (
	"github.com/faiface/beep"
	flag "github.com/spf13/pflag"
	"github.com/viert/go-lame"
	"io"
)

var (
	quality float64
	bitrate int
)

func init() {
	flag.Float64Var(&quality, "lame-quality", 2, "LAME VBR quality")
	flag.IntVar(&bitrate, "lame-bitrate", 160, "LAME output bitrate")
}

type Encoder struct {
	Encoder *lame.Encoder
	Format  beep.Format
}

func (e Encoder) Write(p []byte) (n int, err error) {
	return e.Encoder.Write(p)
}

func (e Encoder) Close() error {
	e.Encoder.Close()
	return nil
}

func NewEncoder(w io.Writer, format beep.Format) (io.WriteCloser, error) {
	var err error
	enc := Encoder{
		Encoder: lame.NewEncoder(w),
		Format:  format,
	}
	if err = enc.Encoder.SetVBR(lame.VBRDefault); err != nil {
		return enc, err
	}
	if err = enc.Encoder.SetVBRQuality(quality); err != nil {
		return enc, err
	}
	return enc, nil
}
