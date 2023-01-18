package encoder

import (
	"context"
	"errors"
	"fmt"
	"github.com/faiface/beep"
	"golang.org/x/time/rate"
	"io"
	"time"
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

	samples := make([][2]float64, format.SampleRate.N(time.Second/10))
	buffer := make([]byte, len(samples)*format.Width())
	limit := rate.NewLimiter(
		rate.Limit(format.SampleRate.N(time.Second)), // Throttle buffer to 1s per 1s
		format.SampleRate.N(20*time.Second),          // Allow burst up to 20s
	)

	for {
		n, ok := s.Stream(samples)
		if !ok {
			return nil
		}

		buf := buffer
		switch format.Precision {
		case 1:
			for _, sample := range samples[:n] {
				buf = buf[format.EncodeUnsigned(buf, sample):]
			}
		case 2, 3:
			for _, sample := range samples[:n] {
				buf = buf[format.EncodeSigned(buf, sample):]
			}
		default:
			return fmt.Errorf("encode: invalid precision: %d", format.Precision)
		}

		if err := limit.WaitN(ctx, n); err != nil {
			if errors.Is(err, context.Canceled) {
				err = nil
			}
			return err
		}

		select {
		case <-ctx.Done():
			return nil
		default:
			_, err := w.Write(buffer[:n*format.Width()])
			if err != nil {
				return err
			}
		}
	}
}
