// Package ring provides a fixed-capacity ring buffer with concurrent-safe
// reads and writes.
package ring

import (
	"iter"
	"sync"
)

// Buffer is a fixed-capacity ring buffer. Push overwrites the oldest element
// once the buffer is full. Reads use a logical index in [0, Len()) where 0
// is the oldest element and Len()-1 is the newest.
type Buffer[T any] struct {
	buf   []T
	head  int // index of the oldest element
	count int
	mu    sync.RWMutex
}

// New creates a buffer with the given capacity.
func New[T any](capacity int) *Buffer[T] {
	return &Buffer[T]{buf: make([]T, capacity)}
}

// Push writes v at the next slot. When the buffer is full, the oldest entry
// is overwritten.
func (r *Buffer[T]) Push(v T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.count < len(r.buf) {
		r.buf[(r.head+r.count)%len(r.buf)] = v
		r.count++
		return
	}
	r.buf[r.head] = v
	r.head = (r.head + 1) % len(r.buf)
}

// At returns the element at logical position i (0 = oldest). Caller must
// ensure 0 <= i < Len().
func (r *Buffer[T]) At(i int) T {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.buf[(r.head+i)%len(r.buf)]
}

// Len returns the number of elements currently held.
func (r *Buffer[T]) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.count
}

// Entries yields each element in chronological order (oldest first).
func (r *Buffer[T]) Entries() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		r.mu.RLock()
		defer r.mu.RUnlock()
		for i := range r.count {
			if !yield(i, r.buf[(r.head+i)%len(r.buf)]) {
				return
			}
		}
	}
}

// Last yields the most recent n elements in chronological order. If the
// buffer holds fewer than n elements, all of them are yielded. The yielded
// index starts at 0 for the first element returned.
func (r *Buffer[T]) Last(n int) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		r.mu.RLock()
		defer r.mu.RUnlock()
		n = min(n, r.count)
		start := r.count - n
		for i := range n {
			if !yield(i, r.buf[(r.head+start+i)%len(r.buf)]) {
				return
			}
		}
	}
}
