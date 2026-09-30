package handlers

import (
	"net/http"

	"gabe565.com/relax-sounds/internal/config"
	"github.com/pocketbase/pocketbase/core"
)

type ConfigResponse struct {
	CastAppID string `json:"castAppId,omitempty"`
}

func Config(conf *config.Config) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		return e.JSON(http.StatusOK, ConfigResponse{CastAppID: conf.CastAppID})
	}
}
