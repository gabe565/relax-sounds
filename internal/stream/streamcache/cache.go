package streamcache

import (
	"sync"
	"time"
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

func (a *Cache) Delete(id string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.Entries, id)
}

func (a *Cache) Close() {
	a.close <- struct{}{}
	a.cleanup(0)
}

func (a *Cache) beginCleanupCron() chan<- struct{} {
	closer := make(chan struct{})
	go func() {
		ticker := time.NewTicker(scanInterval)

		for {
			select {
			case <-closer:
				ticker.Stop()
				return
			case <-ticker.C:
				a.cleanup(cleanAfter)
			}
		}
	}()
	return closer
}

func (a *Cache) cleanup(since time.Duration) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for id, entry := range a.Entries {
		if entry.Mu.TryLock() {
			if time.Since(entry.Accessed) >= since {
				delete(a.Entries, id)
				entry.Mu.Unlock()
				if err := entry.Close(); err != nil {
					entry.Log.Error("Failed to cleanup stream", "error", err)
				}
			} else {
				entry.Mu.Unlock()
			}
		}
	}
}
