package hls

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"gabe565.com/relax-sounds/internal/handlers/mix/stream"
	"gabe565.com/relax-sounds/internal/ring"
	"github.com/pocketbase/pocketbase/core"
)

// MaxSegments is the in-memory ring buffer depth. Older segments are dropped
// as new ones are pushed. Sized comfortably larger than ManifestWindow so the
// client always finds the segments it sees in the manifest.
const MaxSegments = 10

// ManifestWindow is how many segments appear in any single manifest response.
const ManifestWindow = 5

type Entry struct {
	Log     *slog.Logger
	UUID    string
	Preset  string
	Streams stream.Streams

	mu          sync.RWMutex
	buf         *ring.Buffer[*Segment]
	nextSeq     uint64
	closed      bool
	transferred atomic.Uint64

	ready chan struct{} // closed once the producer has buffered bufferAhead of audio

	readyOnce sync.Once

	ctx     context.Context
	cancel  context.CancelFunc
	Created time.Time
}

// Segment is one HLS media segment — a complete, standalone MP3 file
// covering a fixed slice of mixed audio.
type Segment struct {
	Seq      uint64
	Duration time.Duration
	Bytes    []byte
	Created  time.Time
}

func NewEntry(e *core.RequestEvent, uuid, preset string) *Entry {
	ctx, cancel := context.WithCancel(context.Background())
	now := time.Now()
	return &Entry{
		Log: slog.With(
			"userIP", e.RealIP(),
			"userAgent", e.Request.UserAgent(),
			"url", e.Request.URL.String(),
			"id", uuid,
		),
		UUID:    uuid,
		Preset:  preset,
		buf:     ring.New[*Segment](MaxSegments),
		nextSeq: initialSeq(now),
		ready:   make(chan struct{}),
		ctx:     ctx,
		cancel:  cancel,
		Created: now,
	}
}

// initialSeq derives the first sequence number for a new entry from wall
// clock so sequences are monotonic across server restarts, cache evictions,
// and instance failovers. Paired with the producer's absolute pacing, the
// newest sequence at any moment depends only on the time, not on the entry.
func initialSeq(now time.Time) uint64 {
	//nolint:gosec // UnixNano() is positive for any post-1970 time
	return uint64(now.UnixNano() / int64(SegmentDuration()))
}

// seqStart is the wall-clock time at which seq's audio begins.
func seqStart(seq uint64) time.Time {
	//nolint:gosec // seq is derived from a post-1970 UnixNano
	return time.Unix(0, int64(seq)*int64(SegmentDuration()))
}

// Context returns a context that is canceled when the entry is closed.
func (e *Entry) Context() context.Context { return e.ctx }

// PushSegment writes a freshly encoded segment to the ring buffer, assigning
// the next sequence number. It returns the sequence the following segment
// will receive.
func (e *Entry) PushSegment(seg *Segment) uint64 {
	e.mu.Lock()
	defer e.mu.Unlock()
	seg.Seq = e.nextSeq
	e.nextSeq++
	e.buf.Push(seg)
	return e.nextSeq
}

func (e *Entry) markReady() {
	e.readyOnce.Do(func() { close(e.ready) })
}

// GetSegment returns the segment with the given sequence number, or nil if
// it's no longer in the ring buffer (or hasn't been produced yet).
func (e *Entry) GetSegment(seq uint64) *Segment {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if e.buf.Len() == 0 {
		return nil
	}
	//nolint:gosec // buf.Len() is bounded by MaxSegments
	oldest := e.nextSeq - uint64(e.buf.Len())
	if seq < oldest || seq >= e.nextSeq {
		return nil
	}
	//nolint:gosec // (seq - oldest) is bounded by MaxSegments
	return e.buf.At(int(seq - oldest))
}

// WaitReady blocks until the producer has caught up to bufferAhead,
// the caller's context is canceled, or the entry shuts down.
func (e *Entry) WaitReady(ctx context.Context) error {
	select {
	case <-e.ready:
		return nil
	case <-e.ctx.Done():
		return ErrClosed
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (e *Entry) Close() error {
	e.mu.Lock()
	if e.closed {
		e.mu.Unlock()
		return nil
	}
	e.closed = true
	e.cancel()
	e.mu.Unlock()

	return e.Streams.Close()
}

// AddTransferred increments the running total of bytes written to
// clients.
func (e *Entry) AddTransferred(n int64) {
	if n > 0 {
		e.transferred.Add(uint64(n))
	}
}

// TransferredBytes returns the running total of bytes written.
func (e *Entry) TransferredBytes() uint64 {
	return e.transferred.Load()
}

// ErrClosed is returned by operations that fail because the entry is shutting
// down. Producers use this to exit cleanly.
var ErrClosed = errors.New("hls: entry closed")
