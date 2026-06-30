package config

const (
	FlagMaxPresetLen    = "max-preset-length"
	FlagResampleQuality = "resample-quality"
	FlagLAMEQuality     = "lame-quality"

	FlagCacheCleanAfter = "cache-clean-after"

	FlagMetricsEnabled = "metrics-enabled"
	FlagMetricsAddress = "metrics-address"

	FlagDebugEnabled = "debug-enabled"
	FlagDebugAddress = "debug-address"
)

func (c *Config) RegisterFlags() *Config {
	fs := c.App.RootCmd.PersistentFlags()

	fs.IntVar(&c.MaxPresetLen, FlagMaxPresetLen, c.MaxPresetLen, "Maximum number of sounds that a preset can contain")
	fs.IntVar(&c.ResampleQuality, FlagResampleQuality, c.ResampleQuality,
		"Resample quality. Recommend values between 1-4.",
	)
	fs.Float64Var(&c.LAMEQuality, FlagLAMEQuality, c.LAMEQuality, "LAME VBR quality")

	fs.DurationVar(&c.CacheCleanAfter, FlagCacheCleanAfter, c.CacheCleanAfter,
		"How old a cache entry must be before it is cleaned",
	)

	fs.BoolVar(&c.MetricsEnabled, FlagMetricsEnabled, c.MetricsEnabled, "Enables Prometheus metrics API")
	fs.StringVar(&c.MetricsAddress, FlagMetricsAddress, c.MetricsAddress, "Prometheus metrics API listen address")

	fs.BoolVar(&c.DebugEnabled, FlagDebugEnabled, c.DebugEnabled, "Enables debug server")
	fs.StringVar(&c.DebugAddress, FlagDebugAddress, c.DebugAddress, "Debug server listen address")

	return c
}
