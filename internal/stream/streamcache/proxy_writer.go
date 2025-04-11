package streamcache

import (
	"bytes"
	"errors"
	"io"
)

func NewProxyWriter() *ProxyWriter {
	return &ProxyWriter{}
}

type ProxyWriter struct {
	Limit int
	err   error

	w            io.Writer
	buf          bytes.Buffer
	written      int
	totalWritten int64
}

// Write proxies writes to the provided io.Writer up to a limit.
//
// The limit is enforced to prevent overflowing the previously sent Content-Length.
//
// Write always returns a nil error to avoid halting the encoder on failure.
// Check Err to inspect write errors.
//
//nolint:nilerr
func (b *ProxyWriter) Write(p []byte) (int, error) {
	if b.w == nil {
		return len(p), nil
	}

	if err := b.flush(); err != nil {
		return len(p), nil
	}

	if available := b.Limit - b.written; len(p) > available {
		n, err := b.w.Write(p[:available])
		b.written += n
		b.totalWritten += int64(n)
		if err != nil {
			b.err = err
		} else {
			b.err = io.ErrShortWrite
		}
		b.buf.Write(p[n:])
		return len(p), nil
	}

	n, err := b.w.Write(p)
	b.written += n
	b.totalWritten += int64(n)
	if err != nil {
		b.err = err
	} else if b.written >= b.Limit {
		b.err = io.ErrShortWrite
	}

	return len(p), nil
}

func (b *ProxyWriter) flush() error {
	if b.buf.Len() == 0 {
		return nil
	}
	n, err := io.CopyN(b.w, &b.buf, int64(b.Limit-b.written))
	b.written += int(n)
	b.totalWritten += n
	if err != nil && !errors.Is(err, io.EOF) {
		b.err = err
		return err
	}
	return nil
}

func (b *ProxyWriter) Err() error {
	return b.err
}

func (b *ProxyWriter) Buffered() int {
	return b.buf.Len()
}

// SetWriter sets the output writer for the current chunk.
// Buffered overflow data from previous writes is preserved.
func (b *ProxyWriter) SetWriter(w io.Writer) {
	b.w = w
	b.err = nil
	b.written = 0
}

func (b *ProxyWriter) TotalWritten() int64 {
	return b.totalWritten
}

func (b *ProxyWriter) Close() {
	b.buf.Reset()
	b.err = nil
	b.w = nil
	b.written = 0
}
