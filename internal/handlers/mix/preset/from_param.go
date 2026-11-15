package preset

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
)

const maxDecodedSize = 64 << 10

var ErrTooLarge = errors.New("preset is too large")

func FromParam(encoded string) (Preset, error) {
	if base64.RawURLEncoding.DecodedLen(len(encoded)) > maxDecodedSize {
		return Preset{}, ErrTooLarge
	}

	data, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return Preset{}, err
	}

	var r io.ReadCloser
	switch {
	case len(data) != 0 && data[0] == shorthandVersion:
		return fromShorthand(data)
	case bytes.HasPrefix(data, []byte{0x78, 0x9C}):
		r, err = zlib.NewReader(bytes.NewReader(data))
	case bytes.HasPrefix(data, []byte{0x1F, 0x8B}):
		r, err = gzip.NewReader(bytes.NewReader(data))
	default:
		r = io.NopCloser(bytes.NewReader(data))
	}
	if err != nil {
		return Preset{}, err
	}
	defer func() {
		_ = r.Close()
	}()

	lr := &io.LimitedReader{R: r, N: maxDecodedSize + 1}
	var preset Preset
	if err = json.NewDecoder(lr).Decode(&preset); err != nil {
		if lr.N == 0 {
			return Preset{}, ErrTooLarge
		}
		return Preset{}, err
	}

	if err := r.Close(); err != nil {
		return Preset{}, err
	}

	return preset, nil
}
