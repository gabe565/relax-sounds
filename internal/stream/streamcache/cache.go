package streamcache

import (
	"sync"
	"time"
)

type Cache struct {
	entries map[string]*Entry
	close   chan<- struct{}
	mu      sync.Mutex
}

func New() *Cache {
	cache := &Cache{
		entries: make(map[string]*Entry),
	}
	cache.close = cache.beginCleanupCron()

	return cache
}

func (a *Cache) Set(id string, entry *Entry) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	var err error
	if prev, ok := a.entries[id]; ok {
		// Same window is changing streams
		// Destroy old stream then recreate
		err = prev.Close()
	}
	a.entries[id] = entry
	return err
}

func (a *Cache) Get(id string) (*Entry, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	entry, found := a.entries[id]
	if found {
		entry.Accessed = time.Now()
	}
	return entry, found
}

func (a *Cache) Delete(id string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	var err error
	if prev, ok := a.entries[id]; ok {
		// Same window is changing streams
		// Destroy old stream then recreate
		err = prev.Close()
	}
	delete(a.entries, id)
	return err
}

func (a *Cache) Has(id string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	_, ok := a.entries[id]
	return ok
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

	for id, entry := range a.entries {
		if entry.Mu.TryLock() {
			if time.Since(entry.Accessed) >= since {
				delete(a.entries, id)
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
