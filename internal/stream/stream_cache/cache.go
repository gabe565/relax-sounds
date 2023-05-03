package stream_cache

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Cache struct {
	Entries map[string]*Entry
	close   chan<- struct{}
	mu      sync.Mutex
}

func New() *Cache {
	cache := &Cache{
		Entries: make(map[string]*Entry),
	}
	cache.close = cache.beginCleanupCron()

	return cache
}

func (a *Cache) Add(id string, entry *Entry) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Entries[id] = entry
}

func (a *Cache) Get(id string) (*Entry, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if entry, found := a.Entries[id]; found {
		entry.Accessed = time.Now()
		return entry, true
	}
	return nil, false
}

func (a *Cache) Close() {
	a.cleanup(0)
	a.close <- struct{}{}
}

func (a *Cache) beginCleanupCron() chan<- struct{} {
	closer := make(chan struct{})
	go func() {
		ticker := time.NewTicker(15 * time.Minute)

		for {
			select {
			case <-closer:
				ticker.Stop()
				return
			case <-ticker.C:
				a.cleanup(15 * time.Minute)
			}
		}
	}()
	return closer
}

func (a *Cache) cleanup(since time.Duration) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for id, entry := range a.Entries {
		if time.Since(entry.Accessed) >= since {
			log.WithFields(log.Fields{
				"ip":       entry.RemoteAddr,
				"id":       id,
				"accessed": entry.Accessed,
				"age":      time.Since(entry.Created).Truncate(time.Millisecond).String(),
			}).Info("Cleanup stream")
			if err := entry.Close(); err != nil {
				log.WithError(err).WithField("id", id).
					Error("Failed to cleanup stream")
			}
			delete(a.Entries, id)
		}
	}
}
