package streamcache

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuffer_Len(t *testing.T) {
	type fields struct {
		buf *bytes.Buffer
		cap int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{"nil", fields{}, 0},
		{"empty", fields{buf: bytes.NewBuffer(make([]byte, 0, 1))}, 0},
		{"value", fields{buf: bytes.NewBufferString("abc")}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Buffer{
				buf: tt.fields.buf,
				cap: tt.fields.cap,
			}
			assert.Equal(t, tt.want, b.Len())
		})
	}
}

func TestBuffer_Read(t *testing.T) {
	type fields struct {
		buf *bytes.Buffer
		cap int
	}
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantN   int
		want    string
		wantCap int
		wantErr require.ErrorAssertionFunc
	}{
		{"nil", fields{}, args{p: make([]byte, 1)}, 0, "", 0, require.Error},
		{"nil with empty param", fields{}, args{}, 0, "", 0, require.NoError},
		{"empty", fields{buf: bytes.NewBuffer(nil)}, args{p: make([]byte, 1)}, 0, "", 0, require.Error},
		{"value", fields{buf: bytes.NewBufferString("abc")}, args{p: make([]byte, 3)}, 3, "abc", 3, require.NoError},
		{"higher cap", fields{buf: bytes.NewBufferString("abc"), cap: 5}, args{p: make([]byte, 3)}, 3, "abc", 5, require.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Buffer{
				buf: tt.fields.buf,
				cap: tt.fields.cap,
			}
			got, err := b.Read(tt.args.p)
			tt.wantErr(t, err)
			assert.Equal(t, tt.wantN, got)
			assert.Equal(t, tt.want, string(tt.args.p[:got]))
			assert.Equal(t, tt.wantCap, b.cap)
		})
	}
}

func TestBuffer_ReadFrom(t *testing.T) {
	type fields struct {
		buf *bytes.Buffer
		cap int
	}
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantN   int64
		want    string
		wantCap int
		wantErr require.ErrorAssertionFunc
	}{
		{"nil buffer", fields{}, args{r: strings.NewReader("abc")}, 3, "abc", 0, require.NoError},
		{"with buffer", fields{buf: bytes.NewBuffer(nil)}, args{r: strings.NewReader("abc")}, 3, "abc", 0, require.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Buffer{
				buf: tt.fields.buf,
				cap: tt.fields.cap,
			}
			got, err := b.ReadFrom(tt.args.r)
			tt.wantErr(t, err)
			assert.Equal(t, tt.wantN, got)
			assert.Equal(t, tt.want, b.buf.String())
			assert.Equal(t, tt.wantCap, b.cap)
		})
	}
}

func TestBuffer_Write(t *testing.T) {
	type fields struct {
		buf *bytes.Buffer
		cap int
	}
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantN   int
		want    string
		wantErr require.ErrorAssertionFunc
	}{
		{"nil buffer", fields{}, args{p: []byte("abc")}, 3, "abc", require.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Buffer{
				buf: tt.fields.buf,
				cap: tt.fields.cap,
			}
			got, err := b.Write(tt.args.p)
			tt.wantErr(t, err)
			assert.Equal(t, tt.wantN, got)
			assert.Equal(t, tt.want, b.buf.String())
		})
	}
}

func TestBuffer_WriteTo(t *testing.T) {
	type fields struct {
		buf *bytes.Buffer
		cap int
	}
	tests := []struct {
		name    string
		fields  fields
		wantN   int64
		want    string
		wantCap int
		wantErr require.ErrorAssertionFunc
	}{
		{"nil buffer", fields{}, 0, "", 0, require.Error},
		{"with buffer", fields{buf: bytes.NewBufferString("abc")}, 3, "abc", 3, require.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Buffer{
				buf: tt.fields.buf,
				cap: tt.fields.cap,
			}
			w := strings.Builder{}
			got, err := b.WriteTo(&w)
			tt.wantErr(t, err)
			assert.Equal(t, tt.wantN, got)
			assert.Equal(t, tt.want, w.String())
			assert.Equal(t, tt.wantCap, b.cap)
		})
	}
}

func TestBuffer_ensureBuffer(t *testing.T) {
	type fields struct {
		buf *bytes.Buffer
		cap int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"has buffer", fields{buf: bytes.NewBufferString("abc")}, "abc"},
		{"nil buffer", fields{}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Buffer{
				buf: tt.fields.buf,
				cap: tt.fields.cap,
			}
			b.ensureBuffer()
			assert.Equal(t, tt.want, b.buf.String())
		})
	}
}

func TestBuffer_raiseCap(t *testing.T) {
	type fields struct {
		buf *bytes.Buffer
		cap int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{"increases", fields{buf: bytes.NewBufferString("abc")}, 3},
		{"unchanged", fields{buf: bytes.NewBufferString("abc"), cap: 5}, 5},
		{"nil buffer", fields{cap: 5}, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Buffer{
				buf: tt.fields.buf,
				cap: tt.fields.cap,
			}
			b.raiseCap()
			assert.Equal(t, tt.want, b.cap)
		})
	}
}
