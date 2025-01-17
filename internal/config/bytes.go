package config

import "github.com/labstack/gommon/bytes"

type Bytes int64

func (b Bytes) String() string {
	return bytes.Format(int64(b))
}

func (b *Bytes) Set(s string) error {
	v, err := bytes.Parse(s)
	if err != nil {
		return err
	}
	*b = Bytes(v)
	return nil
}

func (b Bytes) Type() string {
	return "string"
}
