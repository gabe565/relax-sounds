package metrics

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

func Flags(cmd *cobra.Command) {
	enabledDefault := true
	if env := os.Getenv("METRICS_ENABLED"); env != "" {
		var err error
		enabledDefault, err = strconv.ParseBool(env)
		if err != nil {
			slog.Warn("Failed to parse METRICS_ENABLED env", "error", err)
		}
	}
	cmd.PersistentFlags().Bool("metrics-enabled", enabledDefault, "Enables Prometheus metrics API")

	addressDefault := ":9090"
	if env := os.Getenv("METRICS_ADDRESS"); env != "" {
		addressDefault = env
	}

	cmd.PersistentFlags().String("metrics-address", addressDefault, "Prometheus metrics API listen address")
}

func Serve(cmd *cobra.Command) {
	enabled, err := cmd.PersistentFlags().GetBool("metrics-enabled")
	if err != nil {
		panic(err)
	}
	if !enabled {
		return
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	addr, err := cmd.PersistentFlags().GetString("metrics-address")
	if err != nil {
		panic(err)
	}

	slog.Info("Starting metrics server", "address", addr)
	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			slog.Error("Failed to serve metrics", "error", err)
		}
	}()
}
