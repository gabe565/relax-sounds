package mixer

import (
	"errors"
	"github.com/gabe565/relax-sounds/internal/playlist"
	"github.com/gabe565/relax-sounds/internal/stream"
	"github.com/viert/go-lame"
	"net/http"
	"os"
	"syscall"
)

func Mix(res http.ResponseWriter, req *http.Request) {
	var err error
	ctx := req.Context()

	// Stream headers
	res.Header().Set("Connection", "Keep-Alive")
	res.Header().Set("Transfer-Encoding", "chunked")
	res.Header().Set("Content-Type", "audio/mp3")

	// Set up stream
	s, err := stream.New(ctx.Value("playlist").(playlist.Playlist))
	defer s.Close()
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
	if err = encoder.SetQuality(9); err != nil {
		panic(err)
	}
	if err = encoder.SetBrate(160); err != nil {
		panic(err)
	}

	// Encode to wav in a Goroutine
	err = Encode(ctx, encoder, s.Mix(), s.Formats[0])
	if err != nil {
		// Ignore broken pipe errors instead of using a context-aware reader
		if !errors.Is(err, syscall.EPIPE) {
			panic(err)
		}
	}
}
