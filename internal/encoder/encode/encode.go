package encode

import (
	"context"
	"errors"
	"fmt"

	"gabe565.com/relax-sounds/internal/stream/streamcache"
)

var (
	ErrInvalidChannels      = errors.New("invalid number of channels")
	ErrUnsupportedPrecision = errors.New("unsupported precision")
)

// Encode writes a duration of the audio stream in PCM format.
//
// Format precision must be 1 or 2 bytes.
func Encode(ctx context.Context, entry *streamcache.Entry) (int, error) {
	if entry.Format.NumChannels <= 0 {
		return 0, fmt.Errorf("%w: %d", ErrInvalidChannels, entry.Format.NumChannels)
	}

	samples := make([][2]float64, 512)
	buffer := make([]byte, len(samples)*entry.Format.Width())
	var written int

	for {
		n, ok := entry.Mix.Stream(samples)
		if !ok {
			return 0, nil
		}

		buf := buffer
		switch entry.Format.Precision {
		case 1:
			for _, sample := range samples[:n] {
				buf = buf[entry.Format.EncodeUnsigned(buf, sample):]
			}
		case 2, 3:
			for _, sample := range samples[:n] {
				buf = buf[entry.Format.EncodeSigned(buf, sample):]
			}
		default:
			return written, fmt.Errorf("%w: %d", ErrUnsupportedPrecision, entry.Format.Precision)
		}

		if ctx.Err() != nil {
			return written, ctx.Err()
		}

		n, err := entry.Encoder.Write(buffer[:n*entry.Format.Width()])
		written += n
		if err != nil {
			return written, err
		}

		if entry.Writer.Err != nil {
			return written, entry.Writer.Err
		}
	}
}
