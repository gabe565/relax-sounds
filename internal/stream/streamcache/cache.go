package streamcache

import (
	"sync"
	"time"

	"gabe565.com/relax-sounds/internal/config"
	"gabe565.com/utils/must"
	"github.com/pocketbase/pocketbase"
)

type Cache struct {
	conf    *config.Config
	entries map[string]*Entry
	close   chan<- struct{}
	mu      sync.Mutex
}

func New(conf *config.Config) *Cache {
	cache := &Cache{
		conf:    conf,
		entries: make(map[string]*Entry),
	}
	must.Must(cache.RegisterCron(conf.App, ""))
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
		entry.Mu.Lock()
		defer entry.Mu.Unlock()
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

func (a *Cache) RegisterCron(app *pocketbase.PocketBase, id string) error {
	return app.Cron().Add("mixStreamCleanup"+id, "* * * * *", func() {
		a.cleanup(a.conf.CacheCleanAfter)
	})
}

func (a *Cache) cleanup(since time.Duration) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for id, entry := range a.entries {
		if entry.Mu.TryLock() {
			if time.Since(entry.Accessed) >= since {
				delete(a.entries, id)
				entry.Mu.Unlock()
				go func() {
					if err := entry.Close(); err != nil {
						entry.Log.Error("Failed to cleanup stream", "error", err)
					}
				}()
			} else {
				entry.Mu.Unlock()
			}
		}
	}
}
