package debug

//nolint:revive,gosec
import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var (
	enabled bool
	addr    string
)

func Flags(cmd *cobra.Command) {
	enabledDefault := false
	if env := os.Getenv("DEBUG_ENABLED"); env != "" {
		var err error
		enabledDefault, err = strconv.ParseBool(env)
		if err != nil {
			log.WithError(err).Warn("Failed to parse DEBUG_ENABLED")
		}
	}
	cmd.PersistentFlags().BoolVar(&enabled, "debug-enabled", enabledDefault, "Enables debug server")

	addressDefault := ":6060"
	if env := os.Getenv("DEBUG_ADDRESS"); env != "" {
		addressDefault = env
	}

	cmd.PersistentFlags().StringVar(&addr, "debug-address", addressDefault, "Debug server listen address")
}

func Serve() error {
	if !enabled {
		return nil
	}

	server := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 3 * time.Second,
	}
	return server.ListenAndServe()
}
