package preset

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/utils/ptr"
)

func TestFromParam(t *testing.T) {
	type args struct {
		encoded string
	}
	tests := []struct {
		name    string
		args    args
		want    Preset
		wantErr require.ErrorAssertionFunc
	}{
		{"empty", args{}, Preset{}, require.Error},
		{
			"valid",
			args{"H4sIAOu4JWcAA4uuVspMUbJSKkktLlHSUSrLzynNTVWyMtRRKkosgTAKEvOUrAxqYwExJrb2KwAAAA"},
			Preset{{ID: "test", Volume: ptr.To(1.0), Rate: ptr.To(1.0), Pan: ptr.To(0.0)}},
			require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromParam(tt.args.encoded)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
