package hls

import (
	"context"
	"time"

	"gabe565.com/relax-sounds/internal/config"
	"github.com/jellydator/ttlcache/v3"
)

func newCache(conf *config.Config) *Cache {
	c := ttlcache.New[string, *Entry](
		ttlcache.WithTTL[string, *Entry](conf.CacheCleanAfter),
	)

	c.OnInsertion(func(_ context.Context, i *ttlcache.Item[string, *Entry]) {
		i.Value().Log.Info("Create HLS stream")
		activeStreamMetrics.Inc()
		totalStreamMetrics.Inc()
	})

	c.OnEviction(func(_ context.Context, _ ttlcache.EvictionReason, i *ttlcache.Item[string, *Entry]) {
		e := i.Value()
		e.Log.Info("Close HLS stream",
			"age", time.Since(e.Created).Round(100*time.Millisecond).String(),
			"transferred", config.Bytes(e.TransferredBytes()).String(), //nolint:gosec // small total
		)
		if err := e.Close(); err != nil {
			e.Log.Error("Failed to cleanup HLS stream", "error", err)
		}
		activeStreamMetrics.Dec()
	})

	go c.Start()

	return &Cache{Cache: c}
}

// Cache wraps ttlcache.Cache so Set on an existing key properly tears down the
// previous entry's producer goroutine and stream handles, instead of silently
// orphaning them (the underlying library updates in place without firing
// OnEviction).
type Cache struct {
	*ttlcache.Cache[string, *Entry]
}

func (c *Cache) Set(key string, value *Entry, ttl time.Duration) {
	c.Delete(key)
	c.Cache.Set(key, value, ttl)
}
