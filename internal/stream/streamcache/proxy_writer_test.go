package streamcache

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProxyWriter(t *testing.T) {
	w := NewProxyWriter()
	w.Limit = 3

	var buf strings.Builder
	w.SetWriter(&buf)
	n, err := w.Write([]byte("hello"))
	require.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, "hel", buf.String())
	require.Error(t, w.err)
	assert.Equal(t, 3, w.written)

	buf.Reset()
	w.SetWriter(&buf)
	require.NoError(t, w.err)
	assert.Zero(t, w.written)
	n, err = w.Write([]byte("world"))
	require.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, "low", buf.String())
	require.Error(t, w.err)
	assert.Equal(t, 3, w.written)

	buf.Reset()
	w.SetWriter(&buf)
	w.Limit = 10
	n, err = w.Write(nil)
	require.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.Equal(t, "orld", buf.String())
	require.NoError(t, w.err)
	assert.Equal(t, 4, w.written)
}
