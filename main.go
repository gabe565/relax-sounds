package main

import (
	"github.com/gabe565/relax-sounds/internal/handlers"
	_ "github.com/gabe565/relax-sounds/migrations"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	app := pocketbase.New()

	migratecmd.MustRegister(app, app.RootCmd, &migratecmd.Options{
		Automigrate: automigrateEnabled(),
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/*", handlers.StaticHandler())
		e.Router.GET("/api/mix/:uuid/:query", handlers.Mix(app))
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
