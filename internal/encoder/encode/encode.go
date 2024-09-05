package encode

import (
	"context"
	"fmt"
	"time"

	"github.com/gabe565/relax-sounds/internal/stream/streamcache"
)

// Encode writes a duration of the audio stream in PCM format.
//
// Format precision must be 1 or 2 bytes.
func Encode(ctx context.Context, entry *streamcache.Entry) error {
	durationPerLoop := time.Second / 10
	samples := make([][2]float64, entry.Format.SampleRate.N(durationPerLoop))
	buffer := make([]byte, len(samples)*entry.Format.Width())

	for {
		n, ok := entry.Mix.Stream(samples)
		if !ok {
			return nil
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
			return fmt.Errorf("%w: %d", ErrUnsupportedPrecision, entry.Format.Precision)
		}

		select {
		case <-ctx.Done():
			return nil
		default:
			if _, err := entry.Encoder.Write(buffer[:n*entry.Format.Width()]); err != nil {
				return err
			}
			if entry.Writer.Err != nil {
				return entry.Writer.Err
			}
		}
	}
}
