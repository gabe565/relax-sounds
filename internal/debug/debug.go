package debug

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

var (
	enabled bool
	addr    string
)

func init() {
	enabledDefault := false
	if env := os.Getenv("DEBUG_ENABLED"); env != "" {
		var err error
		enabledDefault, err = strconv.ParseBool(env)
		if err != nil {
			log.WithError(err).Warn("Failed to parse DEBUG_ENABLED")
		}
	}
	flag.BoolVar(&enabled, "debug-enabled", enabledDefault, "Enables debug server")

	addressDefault := ":6060"
	if env := os.Getenv("DEBUG_ADDRESS"); env != "" {
		addressDefault = env
	}

	flag.StringVar(&addr, "debug-address", addressDefault, "Debug server listen address")
}

func Serve() error {
	if !enabled {
		return nil
	}

	return http.ListenAndServe(addr, nil)
}
