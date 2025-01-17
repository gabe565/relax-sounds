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
	Err   error

	w            io.Writer
	buf          bytes.Buffer
	written      int
	totalWritten int64
}

func (b *ProxyWriter) Write(p []byte) (int, error) {
	if b.w == nil {
		return len(p), nil
	}

	if b.buf.Len() != 0 {
		// Flush buffer
		n, err := io.CopyN(b.w, &b.buf, int64(b.Limit-b.written))
		b.written += int(n)
		b.totalWritten += n
		if err != nil && !errors.Is(err, io.EOF) {
			b.Err = err
			return len(p), nil
		}
	}

	if available := b.Limit - b.written; len(p) > available {
		n, err := b.w.Write(p[:available])
		b.written += n
		b.totalWritten += int64(n)
		if err != nil {
			b.Err = err
		} else {
			b.Err = io.ErrShortWrite
		}
		b.buf.Write(p[n:])
		return len(p), nil
	}

	n, err := b.w.Write(p)
	b.written += n
	b.totalWritten += int64(n)
	if err != nil {
		b.Err = err
	} else if b.written >= b.Limit {
		b.Err = io.ErrShortWrite
	}

	return len(p), nil
}

func (b *ProxyWriter) Buffered() int {
	return b.buf.Len()
}

func (b *ProxyWriter) SetWriter(w io.Writer) {
	b.w = w
	b.Err = nil
	b.written = 0
}

func (b *ProxyWriter) TotalWritten() int64 {
	return b.totalWritten
}

func (b *ProxyWriter) Close() {
	b.buf.Reset()
	b.Err = nil
	b.w = nil
	b.written = 0
}
