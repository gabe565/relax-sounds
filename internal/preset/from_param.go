package preset

import (
	"encoding/base64"
	"encoding/json"
	"github.com/pocketbase/pocketbase/apis"
)

func FromParam(encoded string) (Preset, error) {
	data, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return Preset{}, apis.NewBadRequestError("", nil)
	}

	var entries Shorthand
	if err = json.Unmarshal(data, &entries); err != nil {
		return Preset{}, apis.NewBadRequestError("", nil)
	}

	preset, err := entries.ToPreset()
	if err != nil {
		return Preset{}, apis.NewBadRequestError("", nil)
	}

	return preset, nil
}
