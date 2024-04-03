package metrics

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var (
	enabled bool
	addr    string
)

func Flags(cmd *cobra.Command) {
	enabledDefault := true
	if env := os.Getenv("METRICS_ENABLED"); env != "" {
		var err error
		enabledDefault, err = strconv.ParseBool(env)
		if err != nil {
			log.WithError(err).Warn("Failed to parse METRICS_ENABLED")
		}
	}
	cmd.PersistentFlags().BoolVar(&enabled, "metrics-enabled", enabledDefault, "Enables Prometheus metrics API")

	addressDefault := ":9090"
	if env := os.Getenv("METRICS_ADDRESS"); env != "" {
		addressDefault = env
	}

	cmd.PersistentFlags().StringVar(&addr, "metrics-address", addressDefault, "Prometheus metrics API listen address")
}

func Serve() error {
	if !enabled {
		return nil
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 3 * time.Second,
	}
	return server.ListenAndServe()
}
