package main

import (
	"github.com/gabe565/relax-sounds/internal/debug"
	"github.com/gabe565/relax-sounds/internal/handlers"
	"github.com/gabe565/relax-sounds/internal/metrics"
	_ "github.com/gabe565/relax-sounds/migrations"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	app := pocketbase.New()

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: automigrateEnabled(),
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/*", handlers.StaticHandler())
		e.Router.GET("/api/mix/:uuid/:query", handlers.Mix(app))
		return nil
	})

	go func() {
		if err := metrics.Serve(); err != nil {
			log.Error(err)
		}
	}()

	go func() {
		if err := debug.Serve(); err != nil {
			log.Error(err)
		}
	}()

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
