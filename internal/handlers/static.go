package handlers

import (
	"gabe565.com/relax-sounds/frontend"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func Static() (func(*core.RequestEvent) error, error) {
	fs, err := frontend.FS()
	if err != nil {
		return nil, err
	}

	return apis.Static(fs, true), nil
}
