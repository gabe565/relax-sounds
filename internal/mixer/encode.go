package mixer

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/faiface/beep"
	"io"
)

type header struct {
	RiffMark      [4]byte
	FileSize      int32
	WaveMark      [4]byte
	FmtMark       [4]byte
	FormatSize    int32
	FormatType    int16
	NumChans      int16
	SampleRate    int32
	ByteRate      int32
	BytesPerFrame int16
	BitsPerSample int16
	DataMark      [4]byte
	DataSize      int32
}

// Encode writes all audio streamed from s to w in WAVE format.
//
// Format precision must be 1 or 2 bytes.
func Encode(ctx context.Context, w io.Writer, s beep.Streamer, format beep.Format, includeHeader bool) (err error) {
	if format.NumChannels <= 0 {
		return errors.New("encode: invalid number of channels (less than 1)")
	}
	if format.Precision != 1 && format.Precision != 2 && format.Precision != 3 {
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

	if includeHeader {
		h := header{
			RiffMark:      [4]byte{'R', 'I', 'F', 'F'},
			FileSize:      -1, // finalization
			WaveMark:      [4]byte{'W', 'A', 'V', 'E'},
			FmtMark:       [4]byte{'f', 'm', 't', ' '},
			FormatSize:    16,
			FormatType:    1,
			NumChans:      int16(format.NumChannels),
			SampleRate:    int32(format.SampleRate),
			ByteRate:      int32(int(format.SampleRate) * format.NumChannels * format.Precision),
			BytesPerFrame: int16(format.NumChannels * format.Precision),
			BitsPerSample: int16(format.Precision) * 8,
			DataMark:      [4]byte{'d', 'a', 't', 'a'},
			DataSize:      -1, // finalization
		}
		if err := binary.Write(w, binary.LittleEndian, &h); err != nil {
			return err
		}
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
