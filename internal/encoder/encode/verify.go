package encode

import (
	"errors"
	"fmt"

	"github.com/faiface/beep"
)

var (
	ErrInvalidChannels      = errors.New("invalid number of channels")
	ErrUnsupportedPrecision = errors.New("unsupported precision")
)

func VerifyFormat(format beep.Format) error {
	if format.NumChannels <= 0 {
		return fmt.Errorf("%w: %d", ErrInvalidChannels, format.NumChannels)
	}

	switch format.Precision {
	case 1, 2, 3:
		//
	default:
		return fmt.Errorf("%w: %d", ErrUnsupportedPrecision, format.Precision)
	}

	return nil
}
