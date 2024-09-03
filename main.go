package main

import (
	"log/slog"

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
					slog.Error("Failed to convert sound", "error", err)
				}
			}
		}()
		return nil
	})

	app.OnBeforeServe().Add(func(_ *core.ServeEvent) error {
		metrics.Serve(app.RootCmd)
		debug.Serve(app.RootCmd)
		slog.SetDefault(app.Logger())
		return nil
	})

	if err := app.Start(); err != nil {
		slog.Error("Failed to start app", "error", err)
	}
}
