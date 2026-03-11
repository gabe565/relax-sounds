package cache

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type errWriter struct {
	n   int
	err error
}

func (w errWriter) Write(_ []byte) (int, error) {
	return w.n, w.err
}

var errWrite = errors.New("write failed")

func TestProxyWriter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
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
	})

	t.Run("write error", func(t *testing.T) {
		w := NewProxyWriter()
		w.Limit = 10
		w.SetWriter(errWriter{
			n:   2,
			err: errWrite,
		})

		n, err := w.Write([]byte("hello"))
		require.NoError(t, err)
		assert.Equal(t, 5, n)
		require.ErrorIs(t, w.Err(), errWrite)
		assert.Equal(t, int64(2), w.TotalWritten())
	})

	t.Run("flush error", func(t *testing.T) {
		w := NewProxyWriter()
		w.Limit = 3

		var okWriter strings.Builder
		w.SetWriter(&okWriter)
		n, err := w.Write([]byte("hello"))
		require.NoError(t, err)
		assert.Equal(t, 5, n)
		require.ErrorIs(t, w.Err(), io.ErrShortWrite)
		assert.Equal(t, "hel", okWriter.String())
		assert.Equal(t, 2, w.Buffered())

		w.SetWriter(errWriter{
			n:   0,
			err: errWrite,
		})
		n, err = w.Write([]byte("x"))
		require.NoError(t, err)
		assert.Equal(t, 1, n)
		require.ErrorIs(t, w.Err(), errWrite)
		assert.Zero(t, w.Buffered())
	})
}
