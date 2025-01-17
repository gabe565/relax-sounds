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
func Encode(ctx context.Context, entry *streamcache.Entry) error {
	if entry.Format.NumChannels <= 0 {
		return fmt.Errorf("%w: %d", ErrInvalidChannels, entry.Format.NumChannels)
	}

	var (
		samples = make([][2]float64, 512)
		buffer  = make([]byte, len(samples)*entry.Format.Width())
		encode  func(p []byte, sample [2]float64) (n int)
	)

	switch entry.Format.Precision {
	case 1:
		encode = entry.Format.EncodeUnsigned
	case 2, 3:
		encode = entry.Format.EncodeSigned
	default:
		return fmt.Errorf("%w: %d", ErrUnsupportedPrecision, entry.Format.Precision)
	}

	for {
		n, ok := entry.Mix.Stream(samples)
		if !ok {
			return nil
		}

		buf := buffer
		for _, sample := range samples[:n] {
			buf = buf[encode(buf, sample):]
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		_, err := entry.Encoder.Write(buffer[:n*entry.Format.Width()])
		switch {
		case err != nil:
			return err
		case entry.Writer.Err != nil:
			return entry.Writer.Err
		}
	}
}
