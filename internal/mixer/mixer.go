package mixer

import (
	"context"
	"errors"
	"github.com/faiface/beep"
	"github.com/gabe565/relax-sounds/internal/encoder"
	"github.com/gabe565/relax-sounds/internal/preset"
	"github.com/gabe565/relax-sounds/internal/stream"
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
	"syscall"
)

func Mix(res http.ResponseWriter, req *http.Request) {
	var err error
	ctx := req.Context()

	fileType := encoder.FileType(0)
	err = (&fileType).UnmarshalText([]byte(chi.URLParam(req, "filetype")))
	if err != nil {
		if errors.Is(err, encoder.ErrInvalidFileType) {
			http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		panic(err)
	}

	// Stream headers
	res.Header().Set("Connection", "Keep-Alive")
	res.Header().Set("Transfer-Encoding", "chunked")
	res.Header().Set("Content-Type", fileType.ContentType())

	// Set up stream
	s, err := stream.New(ctx.Value(preset.RequestKey).(preset.Preset))
	defer func(s *stream.Stream) {
		_ = s.Close()
	}(&s)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Invalid file ID returns 404
			http.Error(res, http.StatusText(404), 404)
			return
		} else {
			// Other errors return 500
			http.Error(res, http.StatusText(500), 500)
			panic(err)
		}
	}

	format := beep.Format{
		SampleRate:  44100,
		NumChannels: 2,
		Precision:   2,
	}

	// Get current filetype encoder
	enc, err := fileType.NewEncoder(res, format)
	if err != nil {
		panic(err)
	}
	defer func(encoder encoder.Encoder) {
		_ = encoder.Close()
	}(enc)

	// Write mix to encoder
	err = encoder.Encode(ctx, enc, s.Mix(), format)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, syscall.EPIPE) {
			return
		}
		panic(err)
	}
}
