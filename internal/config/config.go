package config

import (
	"time"

	"gabe565.com/utils/bytefmt"
	"github.com/pocketbase/pocketbase"
)

type Config struct {
	App    *pocketbase.PocketBase
	Public string

	MaxPresetLen    int
	ResampleQuality int
	LAMEQuality     float64
	MixTotalSize    Bytes
	MixChunkSize    Bytes

	CacheCleanAfter time.Duration

	ValkeyEnabled  bool
	ValkeyHost     string
	ValkeyPort     uint16
	ValkeyPassword string
	ValkeyDB       int

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
		MixTotalSize:    1.5 * bytefmt.GiB,
		MixChunkSize:    2 * bytefmt.MiB,

		CacheCleanAfter: 15 * time.Minute,

		ValkeyHost: "localhost",
		ValkeyPort: 6379,

		MetricsEnabled: true,
		MetricsAddress: ":9090",

		DebugAddress: ":6060",
	}
}
