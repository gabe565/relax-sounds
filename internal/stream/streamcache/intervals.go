package streamcache

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var (
	scanInterval time.Duration
	cleanAfter   time.Duration
)

func Flags(cmd *cobra.Command) {
	scanIntervalDefault := time.Minute
	if env := os.Getenv("CACHE_SCAN_INTERVAL"); env != "" {
		var err error
		scanIntervalDefault, err = time.ParseDuration(env)
		if err != nil {
			log.WithError(err).Warn("Failed to parse CACHE_SCAN_INTERVAL")
		}
	}
	cmd.PersistentFlags().DurationVar(&scanInterval, "cache-scan-interval", scanIntervalDefault, "Interval to search stream cache for old entries")

	cleanAfterDefault := 8 * time.Minute
	if env := os.Getenv("CACHE_CLEAN_AFTER"); env != "" {
		var err error
		cleanAfterDefault, err = time.ParseDuration(env)
		if err != nil {
			log.WithError(err).Warn("Failed to parse CACHE_CLEAN_AFTER")
		}
	}

	cmd.PersistentFlags().DurationVar(&cleanAfter, "cache-clean-after", cleanAfterDefault, "How old a cache entry must be before it is cleaned")
}
