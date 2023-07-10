package metrics

import (
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

var (
	enabled bool
	addr    string
)

func init() {
	enabledDefault := true
	if env := os.Getenv("METRICS_ENABLED"); env != "" {
		var err error
		enabledDefault, err = strconv.ParseBool(env)
		if err != nil {
			log.WithError(err).Warn("Failed to parse METRICS_ENABLED")
		}
	}
	flag.BoolVar(&enabled, "metrics-enabled", enabledDefault, "Enables Prometheus metrics API")

	addressDefault := ":9090"
	if env := os.Getenv("METRICS_ADDRESS"); env != "" {
		addressDefault = env
	}

	flag.StringVar(&addr, "metrics-address", addressDefault, "Prometheus metrics API listen address")
}

func Serve() error {
	if !enabled {
		return nil
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(addr, mux)
}
