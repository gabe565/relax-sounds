package debug

//nolint:revive,gosec
import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Flags(cmd *cobra.Command) {
	enabledDefault := false
	if env := os.Getenv("DEBUG_ENABLED"); env != "" {
		var err error
		enabledDefault, err = strconv.ParseBool(env)
		if err != nil {
			log.Warn().Err(err).Msg("failed to parse env DEBUG_ENABLED")
		}
	}
	cmd.PersistentFlags().Bool("debug-enabled", enabledDefault, "Enables debug server")

	addressDefault := ":6060"
	if env := os.Getenv("DEBUG_ADDRESS"); env != "" {
		addressDefault = env
	}

	cmd.PersistentFlags().String("debug-address", addressDefault, "Debug server listen address")
}

func Serve(cmd *cobra.Command) error {
	enabled, err := cmd.PersistentFlags().GetBool("debug-enabled")
	if err != nil {
		panic(err)
	}
	if !enabled {
		return nil
	}

	addr, err := cmd.PersistentFlags().GetString("debug-address")
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 3 * time.Second,
	}
	return server.ListenAndServe()
}
