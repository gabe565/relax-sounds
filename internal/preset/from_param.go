package preset

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"

	"github.com/pocketbase/pocketbase/apis"
)

func FromParam(encoded string) (Preset, error) {
	data, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return Preset{}, apis.NewBadRequestError("", nil)
	}

	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return Preset{}, apis.NewBadRequestError("", nil)
	}
	defer func(r *gzip.Reader) {
		_ = r.Close()
	}(r)

	var preset Preset
	if err = json.NewDecoder(r).Decode(&preset.Tracks); err != nil {
		return Preset{}, apis.NewBadRequestError("", nil)
	}

	if err := r.Close(); err != nil {
		return Preset{}, apis.NewBadRequestError("", nil)
	}

	return preset, nil
}
