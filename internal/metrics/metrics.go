package metrics

import (
	"log/slog"
	"net/http"
	"time"

	"gabe565.com/relax-sounds/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Serve(conf *config.Config) {
	if !conf.MetricsEnabled {
		return
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	slog.Info("Starting metrics server", "address", conf.MetricsAddress)
	server := &http.Server{
		Addr:              conf.MetricsAddress,
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			slog.Error("Failed to serve metrics", "error", err)
		}
	}()
}
