package config

import (
	"time"

	"github.com/pocketbase/pocketbase"
)

type Config struct {
	App    *pocketbase.PocketBase
	Public string

	MaxPresetLen    int
	ResampleQuality int
	LAMEQuality     float64

	CacheScanInterval time.Duration
	CacheCleanAfter   time.Duration

	MetricsEnabled bool
	MetricsAddress string

	DebugEnabled bool
	DebugAddress string
}

func New(app *pocketbase.PocketBase) *Config {
	return &Config{
		App:    app,
		Public: "frontend/dist",

		MaxPresetLen:    20,
		ResampleQuality: 3,
		LAMEQuality:     2,

		CacheScanInterval: time.Minute,
		CacheCleanAfter:   15 * time.Minute,

		MetricsEnabled: true,
		MetricsAddress: ":9090",

		DebugAddress: ":6060",
	}
}
