package preset

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromParamShorthand(t *testing.T) {
	tests := []struct {
		name    string
		encoded string
		want    Preset
		wantErr require.ErrorAssertionFunc
	}{
		{
			"indexed",
			"A5WSJxWSUE-RXZIaRJJZLw",
			Preset{
				{ShortID: new(39), Volume: new(0.21)},
				{ShortID: new(80), Volume: new(0.79)},
				{ShortID: new(93)},
				{ShortID: new(26), Volume: new(0.68)},
				{ShortID: new(89), Volume: new(0.47)},
			},
			require.NoError,
		},
		{
			"all params",
			"A5GUBzJL0M4",
			Preset{{
				ShortID: new(7),
				Volume:  new(0.5),
				Rate:    new(0.75),
				Pan:     new(-0.5),
			}},
			require.NoError,
		},
		{
			// Pan is stored signed rather than offset.
			"full left pan",
			"A5GUB2Rk0Jw",
			Preset{{ShortID: new(7), Volume: new(1.0), Rate: new(1.0), Pan: new(-1.0)}},
			require.NoError,
		},
		{
			// Percent is an absolute scale, so volume above unity just works.
			"volume above 100%",
			"A5GSB8yH",
			Preset{{ShortID: new(7), Volume: new(1.35)}},
			require.NoError,
		},
		{
			// Sounds without a short_id fall back to their record ID.
			"literal record IDs",
			"A5KSr3p6enp6enp6enp6enp6ejKRr2lsY3o5em01OTVwamNkaA",
			Preset{
				{ID: "zzzzzzzzzzzzzzz", Volume: new(0.5)},
				{ID: "ilcz9zm595pjcdh"},
			},
			require.NoError,
		},
		{"index above one byte", "A5GRzQEs", Preset{{ShortID: new(300)}}, require.NoError},
		{"empty", "A5A", Preset{}, require.NoError},
		{
			// More than 15 tracks switches msgpack to a 16-bit array header.
			"array16 header",
			"A9wAEJEBkQKRA5EEkQWRBpEHkQiRCZEKkQuRDJENkQ6RD5EQ",
			Preset{
				{ShortID: new(1)}, {ShortID: new(2)}, {ShortID: new(3)}, {ShortID: new(4)},
				{ShortID: new(5)}, {ShortID: new(6)}, {ShortID: new(7)}, {ShortID: new(8)},
				{ShortID: new(9)}, {ShortID: new(10)}, {ShortID: new(11)}, {ShortID: new(12)},
				{ShortID: new(13)}, {ShortID: new(14)}, {ShortID: new(15)}, {ShortID: new(16)},
			},
			require.NoError,
		},
		{"too many fields", "A5GVAAAAAAA", Preset{}, require.Error},
		{"no fields", "A5GQ", Preset{}, require.Error},
		{"id is nil", "A5GRwA", Preset{}, require.Error},
		{"truncated", "A5GSBw", Preset{}, require.Error},
		{"track is not an array", "A5EH", Preset{}, require.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromParam(tt.encoded)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTrackRef(t *testing.T) {
	assert.Equal(t, "abc", Track{ID: "abc"}.Ref())
	assert.Equal(t, "#42", Track{ShortID: new(42)}.Ref())
}
