package preset

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
)

func FromParam(encoded string) (Preset, error) {
	data, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return Preset{}, err
	}

	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return Preset{}, err
	}
	defer func(r *gzip.Reader) {
		_ = r.Close()
	}(r)

	var preset Preset
	if err = json.NewDecoder(r).Decode(&preset); err != nil {
		return Preset{}, err
	}

	if err := r.Close(); err != nil {
		return Preset{}, err
	}

	return preset, nil
}
