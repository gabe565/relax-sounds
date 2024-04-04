package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/gabe565/relax-sounds/internal/debug"
	"github.com/gabe565/relax-sounds/internal/encoder/mp3"
	"github.com/gabe565/relax-sounds/internal/handlers"
	"github.com/gabe565/relax-sounds/internal/hooks"
	"github.com/gabe565/relax-sounds/internal/metrics"
	"github.com/gabe565/relax-sounds/internal/stream"
	"github.com/gabe565/relax-sounds/internal/stream/streamcache"
	"github.com/gabe565/relax-sounds/migrations"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		NoColor:    color.NoColor,
		TimeFormat: "2006/01/02 15:04:05",
	})

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: automigrateEnabled(),
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/*", handlers.StaticHandler(app))
		e.Router.GET("/api/mix/:uuid/:query", handlers.Mix(app))
		return nil
	})

	convertHook := hooks.Convert(app)
	app.OnModelAfterCreate("sounds").Add(convertHook)
	app.OnModelAfterUpdate("sounds").Add(convertHook)

	app.OnBeforeServe().Add(func(_ *core.ServeEvent) error {
		go func() {
			if migrations.ConvertAfterStart {
				if err := hooks.ConvertAll(app); err != nil {
					log.Err(err).Msg("failed to convert sound")
				}
			}
		}()
		return nil
	})

	go func() {
		if err := metrics.Serve(app.RootCmd); err != nil {
			log.Err(err).Msg("failed to serve metrics")
		}
	}()

	go func() {
		if err := debug.Serve(app.RootCmd); err != nil {
			log.Err(err).Msg("failed to serve pprof")
		}
	}()

	if err := app.Start(); err != nil {
		log.Fatal().Err(err).Msg("failed to serve")
	}
}
