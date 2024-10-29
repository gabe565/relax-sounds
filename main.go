package main

import (
	"log/slog"
	"os"
	"time"

	"gabe565.com/relax-sounds/internal/debug"
	"gabe565.com/relax-sounds/internal/encoder/mp3"
	"gabe565.com/relax-sounds/internal/handlers"
	"gabe565.com/relax-sounds/internal/hooks"
	"gabe565.com/relax-sounds/internal/metrics"
	"gabe565.com/relax-sounds/internal/stream"
	"gabe565.com/relax-sounds/internal/stream/streamcache"
	"gabe565.com/relax-sounds/migrations"
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

	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		e.Router.GET("/{path...}", handlers.StaticHandler(app))
		handlers.NewMix(app).RegisterRoutes(e)
		return e.Next()
	})

	convertHook := hooks.Convert(app)
	app.OnModelAfterCreateSuccess("sounds").BindFunc(convertHook)
	app.OnModelAfterUpdateSuccess("sounds").BindFunc(convertHook)

	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		go func() {
			if migrations.ConvertAfterStart {
				if err := hooks.ConvertAll(app); err != nil {
					slog.Error("Failed to convert sound", "error", err)
				}
			}
		}()

		metrics.Serve(app.RootCmd)
		debug.Serve(app.RootCmd)

		slog.SetDefault(slog.New(slogmulti.Fanout(
			app.Logger().Handler(),
			slog.Default().Handler(),
		)))

		return e.Next()
	})

	if err := app.Start(); err != nil {
		slog.Error("Failed to start app", "error", err)
	}
}
