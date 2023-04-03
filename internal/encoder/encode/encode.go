package encode

import (
	"context"
	"fmt"
	"github.com/gabe565/relax-sounds/internal/stream/stream_cache"
	"time"
)

// Encode writes a duration of the audio stream in PCM format.
//
// Format precision must be 1 or 2 bytes.
func Encode(ctx context.Context, duration time.Duration, entry *stream_cache.Entry) (err error) {
	samples := make([][2]float64, entry.Format.SampleRate.N(time.Second/10))
	buffer := make([]byte, len(samples)*entry.Format.Width())

	var totalDuration time.Duration
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
			_, err := entry.Encoder.Write(buffer[:n*entry.Format.Width()])
			if err != nil {
				return err
			}

			totalDuration += time.Second / 10
			if totalDuration >= duration {
				return nil
			}
		}
	}
}
