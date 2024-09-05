package streamcache

import (
	"bytes"
	"io"
)

func NewBuffer(bufCap int) *Buffer {
	return &Buffer{cap: bufCap}
}

type Buffer struct {
	buf *bytes.Buffer
	cap int
}

func (b *Buffer) ensureBuffer() {
	if b.buf == nil {
		b.buf = bytes.NewBuffer(make([]byte, 0, b.cap))
	}
}

func (b *Buffer) raiseCap() {
	if b.buf != nil {
		if bufLen := b.buf.Len(); bufLen > b.cap {
			b.cap = bufLen
		}
	}
}

func (b *Buffer) Write(p []byte) (int, error) {
	b.ensureBuffer()
	return b.buf.Write(p)
}

func (b *Buffer) ReadFrom(r io.Reader) (int64, error) {
	b.ensureBuffer()
	return b.buf.ReadFrom(r)
}

func (b *Buffer) Read(p []byte) (int, error) {
	if b.buf == nil {
		if len(p) == 0 {
			return 0, nil
		}
		return 0, io.EOF
	}
	b.raiseCap()
	n, err := b.buf.Read(p)
	if b.buf.Len() == 0 {
		b.Reset()
	}
	return n, err
}

func (b *Buffer) WriteTo(w io.Writer) (int64, error) {
	if b.buf == nil {
		return 0, io.EOF
	}
	b.raiseCap()
	n, err := b.buf.WriteTo(w)
	if b.buf.Len() == 0 {
		b.Reset()
	}
	return n, err
}

func (b *Buffer) Len() int {
	if b.buf == nil {
		return 0
	}
	return b.buf.Len()
}

func (b *Buffer) Reset() {
	b.buf = nil
}
