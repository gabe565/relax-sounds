package cache

import (
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"gabe565.com/relax-sounds/internal/config"
	"gabe565.com/utils/must"
)

//nolint:gochecknoglobals
var nextID atomic.Int64

type Cache struct {
	conf    *config.Config
	id      int64
	entries map[string]*Entry
	mu      sync.RWMutex
}

func New(conf *config.Config) *Cache {
	cache := &Cache{
		conf:    conf,
		id:      nextID.Add(1),
		entries: make(map[string]*Entry),
	}
	must.Must(cache.RegisterCron())
	return cache
}

func (a *Cache) Set(id string, entry *Entry) error {
	a.mu.Lock()
	prev, ok := a.entries[id]
	a.entries[id] = entry
	a.mu.Unlock()

	var err error
	if ok {
		if err = prev.Close(); err != nil {
			prev.Log.Error("Failed to close stream", "error", err)
		}
	}
	return err
}

func (a *Cache) Get(id string) (*Entry, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	entry, found := a.entries[id]
	return entry, found
}

func (a *Cache) Delete(id string) error {
	a.mu.Lock()
	prev, ok := a.entries[id]
	delete(a.entries, id)
	a.mu.Unlock()

	var err error
	if ok {
		if err = prev.Close(); err != nil {
			prev.Log.Error("Failed to close stream", "error", err)
		}
	}
	return err
}

func (a *Cache) Has(id string) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	_, ok := a.entries[id]
	return ok
}

func (a *Cache) idString() string {
	var id string
	if a.id != 1 {
		id = strconv.FormatInt(a.id, 10)
	}
	return id
}

func (a *Cache) Close() {
	a.conf.App.Cron().Remove("mixStreamCleanup" + a.idString())
	a.cleanup(0)
}

func (a *Cache) RegisterCron() error {
	return a.conf.App.Cron().Add("mixStreamCleanup"+a.idString(), "* * * * *", func() {
		a.cleanup(a.conf.CacheCleanAfter)
	})
}

func (a *Cache) cleanup(cleanupAge time.Duration) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for id, entry := range a.entries {
		if entry.Mu.TryLock() {
			if time.Since(entry.Accessed) >= cleanupAge {
				delete(a.entries, id)
				go func() {
					if err := entry.Close(); err != nil {
						entry.Log.Error("Failed to cleanup stream", "error", err)
					}
				}()
			}
			entry.Mu.Unlock()
		}
	}
}
