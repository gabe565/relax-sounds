package main

import (
	"github.com/gabe565/relax-sounds/internal/debug"
	"github.com/gabe565/relax-sounds/internal/encoder/mp3"
	"github.com/gabe565/relax-sounds/internal/handlers"
	"github.com/gabe565/relax-sounds/internal/metrics"
	"github.com/gabe565/relax-sounds/internal/stream"
	"github.com/gabe565/relax-sounds/internal/stream/streamcache"
	_ "github.com/gabe565/relax-sounds/migrations"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	app := pocketbase.New()
	stream.Flags(app.RootCmd)
	streamcache.Flags(app.RootCmd)
	handlers.MixFlags(app.RootCmd)
	handlers.StaticFlags(app.RootCmd)
	mp3.Flags(app.RootCmd)
	metrics.Flags(app.RootCmd)
	debug.Flags(app.RootCmd)

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
