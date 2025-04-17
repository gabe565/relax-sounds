package config

const (
	FlagPublic = "public"

	FlagMaxPresetLen    = "max-preset-length"
	FlagResampleQuality = "resample-quality"
	FlagLAMEQuality     = "lame-quality"
	FlagMixTotalSize    = "mix-total-size"
	FlagMixChunkSize    = "mix-chunk-size"

	FlagValkeyEnabled  = "valkey-enabled"
	FlagValkeyHost     = "valkey-host"
	FlagValkeyPort     = "valkey-port"
	FlagValkeyPassword = "valkey-password"
	FlagValkeyDB       = "valkey-db"

	FlagCacheCleanAfter = "cache-clean-after"

	FlagMetricsEnabled = "metrics-enabled"
	FlagMetricsAddress = "metrics-address"

	FlagDebugEnabled = "debug-enabled"
	FlagDebugAddress = "debug-address"
)

func (c *Config) RegisterFlags() *Config {
	fs := c.App.RootCmd.PersistentFlags()

	fs.StringVar(&c.Public, FlagPublic, c.Public, "Public directory")

	fs.IntVar(&c.MaxPresetLen, FlagMaxPresetLen, c.MaxPresetLen, "Maximum number of sounds that a preset can contain")
	fs.IntVar(&c.ResampleQuality, FlagResampleQuality, c.ResampleQuality,
		"Resample quality. Recommend values between 1-4.",
	)
	fs.Float64Var(&c.LAMEQuality, FlagLAMEQuality, c.LAMEQuality, "LAME VBR quality")
	fs.Var(&c.MixTotalSize, FlagMixTotalSize, "Total size of a mix stream")
	fs.Var(&c.MixChunkSize, FlagMixChunkSize, "Size of each HTTP response from the mix endpoint")

	fs.DurationVar(&c.CacheCleanAfter, FlagCacheCleanAfter, c.CacheCleanAfter,
		"How old a cache entry must be before it is cleaned",
	)

	fs.BoolVar(&c.ValkeyEnabled, FlagValkeyEnabled, c.ValkeyEnabled, "Enables Valkey position cache")
	fs.StringVar(&c.ValkeyHost, FlagValkeyHost, c.ValkeyHost, "Valkey host")
	fs.Uint16Var(&c.ValkeyPort, FlagValkeyPort, c.ValkeyPort, "Valkey port")
	fs.StringVar(&c.ValkeyPassword, FlagValkeyPassword, c.ValkeyPassword, "Valkey password")
	fs.IntVar(&c.ValkeyDB, FlagValkeyDB, c.ValkeyDB, "Valkey database")

	fs.BoolVar(&c.MetricsEnabled, FlagMetricsEnabled, c.MetricsEnabled, "Enables Prometheus metrics API")
	fs.StringVar(&c.MetricsAddress, FlagMetricsAddress, c.MetricsAddress, "Prometheus metrics API listen address")

	fs.BoolVar(&c.DebugEnabled, FlagDebugEnabled, c.DebugEnabled, "Enables debug server")
	fs.StringVar(&c.DebugAddress, FlagDebugAddress, c.DebugAddress, "Debug server listen address")

	return c
}
