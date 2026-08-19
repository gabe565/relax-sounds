package preset

import (
	"errors"
	"fmt"

	"github.com/vmihailenco/msgpack/v5"
	"github.com/vmihailenco/msgpack/v5/msgpcode"
)

// Binary preset shorthand (version 3).
//
//	byte 0:  version
//	byte 1+: msgpack array of [id, volume?, rate?, pan?] tracks
const shorthandVersion = 3

var ErrInvalidTrack = errors.New("invalid track")

var _ msgpack.CustomDecoder = (*Track)(nil)

func fromPercent(percent int) float64 {
	return float64(percent) / 100
}

func isInt(c byte) bool {
	switch c {
	case msgpcode.Uint8, msgpcode.Uint16, msgpcode.Uint32, msgpcode.Uint64,
		msgpcode.Int8, msgpcode.Int16, msgpcode.Int32, msgpcode.Int64:
		return true
	default:
		return msgpcode.IsFixedNum(c)
	}
}

func (t *Track) DecodeMsgpack(dec *msgpack.Decoder) error {
	n, err := dec.DecodeArrayLen()
	if err != nil {
		return err
	}
	if n < 1 || n > 4 {
		return fmt.Errorf("%w: got %d fields", ErrInvalidTrack, n)
	}

	code, err := dec.PeekCode()
	if err != nil {
		return err
	}

	switch {
	case msgpcode.IsString(code):
		if t.ID, err = dec.DecodeString(); err != nil {
			return err
		}
	case isInt(code):
		index, err := dec.DecodeInt()
		if err != nil {
			return err
		}
		if index < 0 {
			return fmt.Errorf("%w: negative sound index %d", ErrInvalidTrack, index)
		}

		t.ShortID = &index
	default:
		return fmt.Errorf("%w: unexpected sound reference %#x", ErrInvalidTrack, code)
	}

	for i := range n - 1 {
		percent, err := dec.DecodeInt()
		if err != nil {
			return err
		}

		switch i {
		case 0:
			t.Volume = new(fromPercent(percent))
		case 1:
			t.Rate = new(fromPercent(percent))
		case 2:
			t.Pan = new(fromPercent(percent))
		}
	}

	return nil
}

func fromShorthand(data []byte) (Preset, error) {
	var p Preset
	if err := msgpack.Unmarshal(data[1:], &p); err != nil {
		return Preset{}, err
	}
	if p == nil {
		p = Preset{}
	}
	return p, nil
}
