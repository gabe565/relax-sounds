package hls

import "encoding/binary"

const (
	id3Owner      = "com.apple.streaming.transportStreamTimestamp"
	timestampHz   = 90000
	timestampMask = 1<<33 - 1

	id3HeaderLen   = 10
	id3FrameHdrLen = 10
	id3PayloadLen  = len(id3Owner) + 1 + 8
	id3TagLen      = id3HeaderLen + id3FrameHdrLen + id3PayloadLen
)

func timestampFor(samples uint64) uint64 {
	return (samples * timestampHz / segmentSampleRate) & timestampMask
}

func syncsafe(b []byte, n int) {
	b[0] = byte(n >> 21 & 0x7F)
	b[1] = byte(n >> 14 & 0x7F)
	b[2] = byte(n >> 7 & 0x7F)
	b[3] = byte(n & 0x7F)
}

func id3TimestampTag(ts uint64) []byte {
	b := make([]byte, id3TagLen)

	copy(b, "ID3")
	b[3], b[4], b[5] = 4, 0, 0
	syncsafe(b[6:10], id3FrameHdrLen+id3PayloadLen)

	f := b[id3HeaderLen:]
	copy(f, "PRIV")
	syncsafe(f[4:8], id3PayloadLen)
	f[8], f[9] = 0, 0

	p := f[id3FrameHdrLen:]
	n := copy(p, id3Owner)
	p[n] = 0
	binary.BigEndian.PutUint64(p[n+1:], ts)
	return b
}
