package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	flag "github.com/spf13/pflag"
)

var (
	enabled bool
	addr    string
)

func init() {
	flag.BoolVar(&enabled, "metrics-enabled", true, "Enables Prometheus metrics API")
	flag.StringVar(&addr, "metrics-address", ":9090", "Prometheus metrics API listen address")
}

func Serve() error {
	if !enabled {
		return nil
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(addr, mux)
}
