package preset

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"io"
)

func FromParam(encoded string) (Preset, error) {
	data, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return Preset{}, err
	}

	var r io.ReadCloser
	switch {
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

	var preset Preset
	if err = json.NewDecoder(r).Decode(&preset); err != nil {
		return Preset{}, err
	}

	if err := r.Close(); err != nil {
		return Preset{}, err
	}

	return preset, nil
}
