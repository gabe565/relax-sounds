package main

import (
	"log/slog"
	"os"
	"time"

	"gabe565.com/relax-sounds/internal/config"
	"gabe565.com/relax-sounds/internal/debug"
	"gabe565.com/relax-sounds/internal/handlers"
	"gabe565.com/relax-sounds/internal/handlers/mix"
	"gabe565.com/relax-sounds/internal/hooks"
	"gabe565.com/relax-sounds/internal/metrics"
	_ "gabe565.com/relax-sounds/migrations"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
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
	conf := config.New(app).RegisterFlags()

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: automigrateEnabled(),
	})

	app.OnBootstrap().BindFunc(func(e *core.BootstrapEvent) error {
		if err := conf.Load(app.RootCmd); err != nil {
			return err
		}
		return e.Next()
	})

	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		e.Router.BindFunc(func(e *core.RequestEvent) error {
			if e.Request.URL.Path == "/api/health" {
				return apis.SkipSuccessActivityLog().Func(e)
			}
			return e.Next()
		})

		m, err := mix.NewMix(conf)
		if err != nil {
			return err
		}
		m.RegisterRoutes(e)

		e.Router.GET("/{path...}", handlers.StaticHandler(conf))
		metrics.Serve(conf)
		debug.Serve(conf)

		slog.SetDefault(slog.New(slogmulti.Fanout(
			app.Logger().Handler(),
			slog.Default().Handler(),
		)))

		go func() {
			if err := hooks.ConvertAll(app); err != nil {
				slog.Error("Failed to convert sounds", "error", err)
			}
		}()

		return e.Next()
	})

	convertHook := hooks.Convert(app)
	app.OnModelAfterCreateSuccess("sounds").BindFunc(convertHook)
	app.OnModelAfterUpdateSuccess("sounds").BindFunc(convertHook)

	if err := app.Start(); err != nil {
		slog.Error("Failed to start app", "error", err)
		os.Exit(1)
	}
}
