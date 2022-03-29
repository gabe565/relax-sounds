package encoder

import (
	"context"
	"errors"
	"fmt"
	"github.com/faiface/beep"
	"io"
)

// Encode writes all audio streamed from s to w in WAVE format.
//
// Format precision must be 1 or 2 bytes.
func Encode(ctx context.Context, w io.Writer, s beep.Streamer, format beep.Format) (err error) {
	if format.NumChannels <= 0 {
		return errors.New("encode: invalid number of channels (less than 1)")
	}

	switch format.Precision {
	case 1, 2, 3:
		//
	default:
		return errors.New("encode: unsupported precision, 1, 2 or 3 is supported")
	}

	samples := make([][2]float64, 512)
	buffer := make([]byte, len(samples)*format.Width())

	var encode func([]byte, [2]float64) int
	switch format.Precision {
	case 1:
		encode = format.EncodeUnsigned
	case 2, 3:
		encode = format.EncodeSigned
	default:
		return fmt.Errorf("encode: invalid precision: %d", format.Precision)
	}

	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		n, ok := s.Stream(samples)
		if !ok {
			return nil
		}

		buf := buffer
		for _, sample := range samples[:n] {
			buf = buf[encode(buf, sample):]
		}

		_, err := w.Write(buffer[:n*format.Width()])
		if err != nil {
			return err
		}
	}
}
