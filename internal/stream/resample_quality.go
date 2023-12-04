package stream

import (
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

var resampleQuality int

func init() {
	resampleQualityDefault := 3
	if env := os.Getenv("RESAMPLE_QUALITY"); env != "" {
		var err error
		resampleQualityDefault, err = strconv.Atoi(env)
		if err != nil {
			log.WithError(err).Warn("Failed to parse RESAMPLE_QUALITY")
		}
	}
	flag.IntVar(&resampleQuality, "resample-quality", resampleQualityDefault, "Resample quality. Recommend values between 1-4.")
}
