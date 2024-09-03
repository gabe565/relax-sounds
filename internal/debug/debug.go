package debug

//nolint:revive,gosec
import (
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

func Flags(cmd *cobra.Command) {
	enabledDefault := false
	if env := os.Getenv("DEBUG_ENABLED"); env != "" {
		var err error
		enabledDefault, err = strconv.ParseBool(env)
		if err != nil {
			slog.Warn("Failed to parse DEBUG_ENABLED env")
		}
	}
	cmd.PersistentFlags().Bool("debug-enabled", enabledDefault, "Enables debug server")

	addressDefault := ":6060"
	if env := os.Getenv("DEBUG_ADDRESS"); env != "" {
		addressDefault = env
	}

	cmd.PersistentFlags().String("debug-address", addressDefault, "Debug server listen address")
}

func Serve(cmd *cobra.Command) {
	enabled, err := cmd.PersistentFlags().GetBool("debug-enabled")
	if err != nil {
		panic(err)
	}
	if !enabled {
		return
	}

	addr, err := cmd.PersistentFlags().GetString("debug-address")
	if err != nil {
		panic(err)
	}

	slog.Info("Starting debug server", "address", addr)
	server := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 3 * time.Second,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			slog.Error("Failed to serve pprof", "error", err)
		}
	}()
}
