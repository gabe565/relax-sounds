package hls

import (
	"context"
	"testing"
	"time"

	"gabe565.com/relax-sounds/internal/ring"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestEntry(t *testing.T) *Entry {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	return &Entry{
		buf:     ring.New[*Segment](MaxSegments),
		nextSeq: initialSeq(time.Now()),
		ready:   make(chan struct{}),
		pushed:  make(chan struct{}),
		ctx:     ctx,
		cancel:  cancel,
	}
}

func TestReadyAfterStartupSegments(t *testing.T) {
	e := newTestEntry(t)
	for range startupSegments - 1 {
		e.PushSegment(&Segment{})
	}
	select {
	case <-e.ready:
		t.Fatal("ready before startupSegments were produced")
	default:
	}
	e.PushSegment(&Segment{})
	require.NoError(t, e.WaitReady(t.Context()))
}

func TestWaitSegment(t *testing.T) {
	t.Run("existing", func(t *testing.T) {
		e := newTestEntry(t)
		seq := e.NextSeq()
		e.PushSegment(&Segment{})
		seg, err := e.WaitSegment(t.Context(), seq)
		require.NoError(t, err)
		assert.Equal(t, seq, seg.Seq)
	})

	t.Run("waits for producer", func(t *testing.T) {
		e := newTestEntry(t)
		seq := e.NextSeq() + 2
		go func() {
			for range 3 {
				time.Sleep(10 * time.Millisecond)
				e.PushSegment(&Segment{})
			}
		}()
		seg, err := e.WaitSegment(t.Context(), seq)
		require.NoError(t, err)
		assert.Equal(t, seq, seg.Seq)
	})

	t.Run("evicted", func(t *testing.T) {
		e := newTestEntry(t)
		seq := e.NextSeq()
		for range MaxSegments + 1 {
			e.PushSegment(&Segment{})
		}
		_, err := e.WaitSegment(t.Context(), seq)
		require.ErrorIs(t, err, ErrSegmentUnavailable)
	})

	t.Run("too far ahead", func(t *testing.T) {
		e := newTestEntry(t)
		seq := initialSeq(time.Now().Add(bufferAhead + 2*SegmentDuration()))
		_, err := e.WaitSegment(t.Context(), seq)
		require.ErrorIs(t, err, ErrSegmentUnavailable)
	})

	t.Run("entry closed", func(t *testing.T) {
		e := newTestEntry(t)
		seq := e.NextSeq()
		e.cancel()
		_, err := e.WaitSegment(t.Context(), seq)
		require.ErrorIs(t, err, ErrClosed)
	})

	t.Run("request canceled", func(t *testing.T) {
		e := newTestEntry(t)
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		_, err := e.WaitSegment(ctx, e.NextSeq())
		require.ErrorIs(t, err, context.Canceled)
	})
}

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

func TestBufferSizing(t *testing.T) {
	window := ManifestWindow * SegmentDuration()
	assert.GreaterOrEqual(t, bufferAhead, startOffset, "first playlist must already reach startOffset")
	assert.GreaterOrEqual(t, window, startOffset, "playlist must list startOffset of audio")
	assert.GreaterOrEqual(t, MaxSegments, ManifestWindow, "ring must hold every listed segment")
	assert.Greater(t, MaxSegments*SegmentDuration(), bufferAhead, "ring must hold the startup burst")
	assert.LessOrEqual(
		t,
		startupSegments*SegmentDuration(),
		bufferAhead,
		"entry must become ready before the burst ends",
	)
}
