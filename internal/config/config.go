package config

import (
	"time"

	"github.com/pocketbase/pocketbase"
)

type Config struct {
	App *pocketbase.PocketBase

	MaxPresetLen    int
	ResampleQuality int
	LAMEQuality     float64

	CacheCleanAfter time.Duration

	MetricsEnabled bool
	MetricsAddress string

	DebugEnabled bool
	DebugAddress string
}

func New(app *pocketbase.PocketBase) *Config {
	return &Config{
		App: app,

		MaxPresetLen:    20,
		ResampleQuality: 3,
		LAMEQuality:     2,

		CacheCleanAfter: 2 * time.Minute,

		MetricsEnabled: true,
		MetricsAddress: ":9090",

		DebugAddress: ":6060",
	}
}
