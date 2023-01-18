package handlers

import (
	"context"
	"errors"
	"github.com/faiface/beep"
	"github.com/gabe565/relax-sounds/internal/encoder"
	"github.com/gabe565/relax-sounds/internal/encoder/filetype"
	"github.com/gabe565/relax-sounds/internal/preset"
	"github.com/gabe565/relax-sounds/internal/stream"
	"net/http"
	"os"
	"syscall"
)

var globalMixCtx, globalMixCancel = context.WithCancel(context.Background())

func MixCancelFunc() context.CancelFunc {
	return globalMixCancel
}

func Mix(res http.ResponseWriter, req *http.Request) {
	var err error
	ctx := req.Context()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	fileType := ctx.Value(filetype.RequestKey).(filetype.FileType)

	// Stream headers
	res.Header().Set("Connection", "Keep-Alive")
	res.Header().Set("Transfer-Encoding", "chunked")
	res.Header().Set("Content-Type", fileType.ContentType())

	// Set up stream
	s, err := stream.New(ctx.Value(preset.RequestKey).(preset.Preset))
	defer func(s *stream.Streams) {
		_ = s.Close()
	}(&s)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Invalid file ID returns 404
			http.Error(res, http.StatusText(404), 404)
			return
		}
		// Other errors return 500
		panic(err)
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

	go func() {
		select {
		case <-globalMixCtx.Done():
			// Stop streaming when global cancel called
			cancel()
		case <-ctx.Done():
			// Exit Goroutine when request ends
		}
	}()

	// Write mix to encoder
	if err = encoder.Encode(ctx, enc, s.Mix(), format); err != nil {
		if !errors.Is(err, syscall.EPIPE) && !errors.Is(err, syscall.ECONNRESET) {
			panic(err)
		}
	}
}
