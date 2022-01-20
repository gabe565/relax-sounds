package mixer

import (
	"errors"
	"github.com/faiface/beep"
	"github.com/gabe565/relax-sounds/internal/preset"
	"github.com/gabe565/relax-sounds/internal/stream"
	flag "github.com/spf13/pflag"
	"github.com/viert/go-lame"
	"net/http"
	"os"
	"syscall"
)

var (
	quality float64
	bitrate int
)

func init() {
	flag.Float64Var(&quality, "quality", 2, "LAME VBR quality")
	flag.IntVar(&bitrate, "bitrate", 160, "LAME output bitrate")
}

func Mix(res http.ResponseWriter, req *http.Request) {
	var err error
	ctx := req.Context()

	// Stream headers
	res.Header().Set("Connection", "Keep-Alive")
	res.Header().Set("Transfer-Encoding", "chunked")
	res.Header().Set("Content-Type", "audio/mp3")

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

	// Encode wav to mp3
	encoder := lame.NewEncoder(res)
	defer encoder.Close()
	if err = encoder.SetVBR(lame.VBRDefault); err != nil {
		panic(err)
	}
	if err = encoder.SetVBRQuality(quality); err != nil {
		panic(err)
	}

	// Encode to wav in a Goroutine
	err = Encode(ctx, encoder, s.Mix(), beep.Format{
		SampleRate:  44100,
		NumChannels: 2,
		Precision:   2,
	}, false)
	if err != nil {
		// Ignore broken pipe errors instead of using a context-aware reader
		if !errors.Is(err, syscall.EPIPE) {
			panic(err)
		}
	}
}
