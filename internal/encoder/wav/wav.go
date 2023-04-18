package wav

import (
	"encoding/binary"
	"io"

	"github.com/faiface/beep"
)

type Encoder struct {
	w      io.Writer
	Format beep.Format
}

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

func (e Encoder) Write(p []byte) (n int, err error) {
	return e.w.Write(p)
}

func (e Encoder) Close() error {
	return nil
}

func NewEncoder(w io.Writer, format beep.Format) (io.WriteCloser, error) {
	encoder := Encoder{
		w:      w,
		Format: format,
	}

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
		return encoder, err
	}

	return encoder, nil
}
