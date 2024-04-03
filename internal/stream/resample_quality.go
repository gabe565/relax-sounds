package stream

import (
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
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
			log.Warn().Err(err).Msg("failed to parse env RESAMPLE_QUALITY")
		}
	}
	cmd.PersistentFlags().IntVar(&resampleQuality, "resample-quality", resampleQualityDefault, "Resample quality. Recommend values between 1-4.")
}
