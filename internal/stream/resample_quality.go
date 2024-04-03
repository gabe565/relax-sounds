package stream

import (
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var resampleQuality int

func Flags(cmd *cobra.Command) {
	resampleQualityDefault := 3
	if env := os.Getenv("RESAMPLE_QUALITY"); env != "" {
		var err error
		resampleQualityDefault, err = strconv.Atoi(env)
		if err != nil {
			log.WithError(err).Warn("Failed to parse RESAMPLE_QUALITY")
		}
	}
	cmd.PersistentFlags().IntVar(&resampleQuality, "resample-quality", resampleQualityDefault, "Resample quality. Recommend values between 1-4.")
}
