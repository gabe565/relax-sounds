package preset

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			"zlib",
			args{"eJyKrlbKTFGyUipJLS5R0lEqy88pzU1VsjLUUSpKLIEwChLzlKwMamMBAQAA__8lzA0T"},
			Preset{{ID: "test", Volume: new(1.0), Rate: new(1.0), Pan: new(0.0)}},
			require.NoError,
		},
		{
			"gzip",
			args{"H4sIAOu4JWcAA4uuVspMUbJSKkktLlHSUSrLzynNTVWyMtRRKkosgTAKEvOUrAxqYwExJrb2KwAAAA"},
			Preset{{ID: "test", Volume: new(1.0), Rate: new(1.0), Pan: new(0.0)}},
			require.NoError,
		},
		{
			"plain",
			args{"W3siaWQiOiJ0ZXN0Iiwidm9sdW1lIjoxLCJyYXRlIjoxLCJwYW4iOjB9XQ"},
			Preset{{ID: "test", Volume: new(1.0), Rate: new(1.0), Pan: new(0.0)}},
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

func TestFromParamTooLarge(t *testing.T) {
	t.Run("encoded", func(t *testing.T) {
		encoded := base64.RawURLEncoding.EncodeToString(bytes.Repeat([]byte{' '}, maxDecodedSize+1))
		_, err := FromParam(encoded)
		require.ErrorIs(t, err, ErrTooLarge)
	})

	t.Run("decompressed", func(t *testing.T) {
		var buf bytes.Buffer
		w := zlib.NewWriter(&buf)
		_, err := w.Write([]byte("["))
		require.NoError(t, err)
		_, err = w.Write(bytes.Repeat([]byte("{},"), maxDecodedSize))
		require.NoError(t, err)
		_, err = w.Write([]byte("{}]"))
		require.NoError(t, err)
		require.NoError(t, w.Close())

		_, err = FromParam(base64.RawURLEncoding.EncodeToString(buf.Bytes()))
		require.ErrorIs(t, err, ErrTooLarge)
	})
}
