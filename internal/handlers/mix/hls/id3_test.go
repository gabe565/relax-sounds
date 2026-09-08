package hls

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestID3TimestampTag(t *testing.T) {
	const ts = 0x1_2345_6789
	b := id3TimestampTag(ts)

	require.Len(t, b, id3TagLen)
	assert.Equal(t, "ID3", string(b[:3]))
	assert.Equal(t, []byte{4, 0, 0}, b[3:6])
	assert.Equal(t, id3TagLen-id3HeaderLen, syncsafeDecode(t, b[6:10]))

	f := b[id3HeaderLen:]
	assert.Equal(t, "PRIV", string(f[:4]))
	assert.Equal(t, id3PayloadLen, syncsafeDecode(t, f[4:8]))
	assert.Equal(t, []byte{0, 0}, f[8:10])

	p := f[id3FrameHdrLen:]
	assert.Equal(t, id3Owner, string(p[:len(id3Owner)]))
	assert.Zero(t, p[len(id3Owner)])
	assert.Equal(t, uint64(ts), binary.BigEndian.Uint64(p[len(id3Owner)+1:]))
}

func TestSyncsafe(t *testing.T) {
	var b [4]byte
	syncsafe(b[:], 0x0FFFFFFF)
	assert.Equal(t, []byte{0x7F, 0x7F, 0x7F, 0x7F}, b[:])

	syncsafe(b[:], 128)
	assert.Equal(t, []byte{0, 0, 1, 0}, b[:])
}

func TestTimestampFor(t *testing.T) {
	assert.Equal(t, uint64(0), timestampFor(0))
	assert.Equal(t, uint64(timestampHz), timestampFor(segmentSampleRate))

	const samples = segmentSampleRate * 100_000
	const ticks = samples * timestampHz / segmentSampleRate
	require.Greater(t, uint64(ticks), uint64(timestampMask))
	assert.Equal(t, uint64(ticks%(timestampMask+1)), timestampFor(samples))
}

func syncsafeDecode(t *testing.T, b []byte) int {
	t.Helper()
	for i, v := range b {
		require.Zero(t, v&0x80, "byte %d not synchsafe", i)
	}
	return int(b[0])<<21 | int(b[1])<<14 | int(b[2])<<7 | int(b[3])
}
