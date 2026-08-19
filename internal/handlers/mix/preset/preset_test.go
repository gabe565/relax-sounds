package preset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrackGetters(t *testing.T) {
	tests := []struct {
		name              string
		track             Track
		volume, rate, pan float64
	}{
		{"defaults", Track{}, 1, 1, 0},
		{"in range", Track{Volume: new(0.5), Rate: new(0.75), Pan: new(-0.5)}, 0.5, 0.75, -0.5},
		{"boosted volume", Track{Volume: new(1.25)}, 1, 1, 0},
		{"clamps high", Track{Volume: new(99.0), Rate: new(9.0), Pan: new(9.0)}, 1, 1.5, 1},
		{"clamps low", Track{Volume: new(-9.0), Rate: new(-9.0), Pan: new(-9.0)}, 0, 0.5, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.volume, tt.track.GetVolume(), 1e-9)
			assert.InDelta(t, tt.rate, tt.track.GetRate(), 1e-9)
			assert.InDelta(t, tt.pan, tt.track.GetPan(), 1e-9)
		})
	}
}
