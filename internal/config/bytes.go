package config

import (
	"gabe565.com/utils/bytefmt"
)

//nolint:recvcheck
type Bytes int64

func (b Bytes) String() string {
	return bytefmt.Encode(int64(b))
}

func (b *Bytes) Set(s string) error {
	v, err := bytefmt.Decode(s)
	if err != nil {
		return err
	}
	*b = Bytes(v)
	return nil
}

func (b Bytes) Type() string {
	return "string"
}
