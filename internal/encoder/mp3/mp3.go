package mp3

import (
	"io"

	"github.com/gopxl/beep/v2"
	"github.com/spf13/cobra"
	"github.com/viert/go-lame"
)

//nolint:gochecknoglobals
var quality float64

func Flags(cmd *cobra.Command) {
	cmd.PersistentFlags().Float64Var(&quality, "lame-quality", 2, "LAME VBR quality")
}

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
