package encoder

import "io"

type Encoder interface {
	io.WriteCloser
}
