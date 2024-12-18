package handlers

import (
	"os"

	"gabe565.com/relax-sounds/internal/config"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func StaticHandler(conf *config.Config) func(*core.RequestEvent) error {
	return apis.Static(os.DirFS(conf.Public), true)
}
