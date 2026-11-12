package config

import (
	"time"

	"github.com/pocketbase/pocketbase"
)

type Config struct {
	App *pocketbase.PocketBase

	MaxPresetLen    int
	MaxStreamsPerIP int
	ResampleQuality int
	LAMEQuality     float64

	CacheCleanAfter time.Duration

	CastAppID string

	MetricsEnabled bool
	MetricsAddress string

	DebugEnabled bool
	DebugAddress string
}

func New(app *pocketbase.PocketBase) *Config {
	return &Config{
		App: app,

		MaxPresetLen:    20,
		MaxStreamsPerIP: 5,
		ResampleQuality: 3,
		LAMEQuality:     2,

		CacheCleanAfter: 10 * time.Minute,

		CastAppID: "D7AAD9CA",

		MetricsEnabled: true,
		MetricsAddress: ":9090",

		DebugAddress: ":6060",
	}
}
