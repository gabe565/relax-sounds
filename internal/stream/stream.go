package stream

import (
	"github.com/gabe565/relax-sounds/internal/preset"
	"golang.org/x/sync/errgroup"
	"sync"
)

func New(p preset.Preset) (stream Streams, err error) {
	s := make(Streams, 0, len(p.Tracks))

	var mu sync.Mutex
	group := errgroup.Group{}

	for _, entry := range p.Tracks {
		entry := entry
		group.Go(func() error {
			return s.Add(p.Dir, entry, &mu)
		})
	}

	err = group.Wait()
	return s, err
}
