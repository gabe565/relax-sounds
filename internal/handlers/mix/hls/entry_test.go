package hls

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInitialSeq(t *testing.T) {
	now := time.Date(2026, 9, 24, 1, 35, 53, 123456789, time.UTC)
	seq := initialSeq(now)

	assert.False(t, seqStart(seq).After(now))
	assert.True(t, seqStart(seq+1).After(now))
	assert.Equal(t, SegmentDuration(), seqStart(seq+1).Sub(seqStart(seq)))
}

func TestInitialSeqMonotonic(t *testing.T) {
	now := time.Date(2026, 9, 24, 1, 35, 53, 0, time.UTC)
	assert.Equal(t, initialSeq(now)+1, initialSeq(now.Add(SegmentDuration())))
	assert.LessOrEqual(t, initialSeq(now), initialSeq(now.Add(time.Nanosecond)))
}
