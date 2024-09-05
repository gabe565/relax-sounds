package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/gabe565/relax-sounds/internal/debug"
	"github.com/gabe565/relax-sounds/internal/encoder/mp3"
	"github.com/gabe565/relax-sounds/internal/handlers"
	"github.com/gabe565/relax-sounds/internal/hooks"
	"github.com/gabe565/relax-sounds/internal/metrics"
	"github.com/gabe565/relax-sounds/internal/stream"
	"github.com/gabe565/relax-sounds/internal/stream/streamcache"
	"github.com/gabe565/relax-sounds/migrations"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	slogmulti "github.com/samber/slog-multi"
)

func main() {
	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		Level:      slog.LevelInfo,
		TimeFormat: time.DateTime,
		NoColor:    !isatty.IsTerminal(os.Stderr.Fd()) && !isatty.IsCygwinTerminal(os.Stderr.Fd()),
	})))

	app := pocketbase.New()
	stream.Flags(app.RootCmd)
	streamcache.Flags(app.RootCmd)
	handlers.StaticFlags(app.RootCmd)
	mp3.Flags(app.RootCmd)
	metrics.Flags(app.RootCmd)
	debug.Flags(app.RootCmd)

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: automigrateEnabled(),
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/*", handlers.StaticHandler(app))
		handlers.NewMix(app).RegisterRoutes(e.Router)
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

		slog.SetDefault(slog.New(slogmulti.Fanout(
			app.Logger().Handler(),
			slog.Default().Handler(),
		)))

		return nil
	})

	if err := app.Start(); err != nil {
		slog.Error("Failed to start app", "error", err)
	}
}
