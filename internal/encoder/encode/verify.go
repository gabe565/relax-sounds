package encode

import (
	"errors"
	"github.com/faiface/beep"
)

func VerifyFormat(format beep.Format) error {
	if format.NumChannels <= 0 {
		return errors.New("encode: invalid number of channels (less than 1)")
	}

	switch format.Precision {
	case 1, 2, 3:
		//
	default:
		return errors.New("encode: unsupported precision, 1, 2 or 3 is supported")
	}

	return nil
}
