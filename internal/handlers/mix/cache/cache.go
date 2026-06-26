package cache

import (
	"context"
	"time"

	"gabe565.com/relax-sounds/internal/config"
	"github.com/jellydator/ttlcache/v3"
)

func New(conf *config.Config) *Cache {
	cache := ttlcache.New[string, *Entry](
		ttlcache.WithTTL[string, *Entry](conf.CacheCleanAfter),
	)

	cache.OnInsertion(func(_ context.Context, i *ttlcache.Item[string, *Entry]) {
		i.Value().Log.Info("Create stream")

		activeStreamMetrics.Inc()
		totalStreamMetrics.Inc()
	})

	cache.OnEviction(func(_ context.Context, _ ttlcache.EvictionReason, i *ttlcache.Item[string, *Entry]) {
		e := i.Value()

		e.Mu.Lock()
		defer e.Mu.Unlock()

		e.Log.Info("Close stream",
			"age", time.Since(e.Created).Round(100*time.Millisecond).String(),
			"transferred", config.Bytes(e.Writer.TotalWritten()).String(),
		)

		if err := e.close(); err != nil {
			e.Log.Error("Failed to cleanup stream", "error", err)
		}

		activeStreamMetrics.Dec()
	})

	go cache.Start()

	return &Cache{Cache: cache}
}

type Cache struct {
	*ttlcache.Cache[string, *Entry]
}

func (c *Cache) Set(key string, value *Entry, ttl time.Duration) {
	c.Delete(key)
	c.Cache.Set(key, value, ttl)
}
