package metrics

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Flags(cmd *cobra.Command) {
	enabledDefault := true
	if env := os.Getenv("METRICS_ENABLED"); env != "" {
		var err error
		enabledDefault, err = strconv.ParseBool(env)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to parse env METRICS_ENABLED")
		}
	}
	cmd.PersistentFlags().Bool("metrics-enabled", enabledDefault, "Enables Prometheus metrics API")

	addressDefault := ":9090"
	if env := os.Getenv("METRICS_ADDRESS"); env != "" {
		addressDefault = env
	}

	cmd.PersistentFlags().String("metrics-address", addressDefault, "Prometheus metrics API listen address")
}

func Serve(cmd *cobra.Command) error {
	enabled, err := cmd.PersistentFlags().GetBool("metrics-enabled")
	if err != nil {
		panic(err)
	}
	if !enabled {
		return nil
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	addr, err := cmd.PersistentFlags().GetString("metrics-address")
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 3 * time.Second,
	}
	return server.ListenAndServe()
}
