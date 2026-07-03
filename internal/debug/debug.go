package debug

import (
	"log/slog"
	"net/http"
	_ "net/http/pprof" //nolint:gosec
	"time"

	"gabe565.com/relax-sounds/internal/config"
)

func Serve(conf *config.Config) {
	if !conf.DebugEnabled {
		return
	}

	slog.Info("Starting debug server", "address", conf.DebugAddress)
	server := &http.Server{
		Addr:              conf.DebugAddress,
		ReadHeaderTimeout: 3 * time.Second,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			slog.Error("Failed to serve pprof", "error", err)
		}
	}()
}
