package cache

import (
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
)

type CacheItem struct {
	Value     string
	ExpiresAt time.Time
}

type SnapshotItem struct {
	Key     string
	Value   string
	TTLLeft time.Duration
}

type ChronoCache struct {
	cache      *lru.Cache[string, CacheItem]
	lock       sync.Mutex
	defaultTTL time.Duration
}

func NewChronoCache(size int, defaultTTL time.Duration) (*ChronoCache, error) {
	l, err := lru.New[string, CacheItem](size)
	if err != nil {
		return nil, err
	}
	return &ChronoCache{
		cache:      l,
		defaultTTL: defaultTTL,
	}, nil
}

func (c *ChronoCache) SetWithTTL(key, value string, ttl time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cache.Add(key, CacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	})
}

func (c *ChronoCache) Get(key string) (string, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	item, ok := c.cache.Get(key)
	if !ok || time.Now().After(item.ExpiresAt) {
		c.cache.Remove(key)
		return "", false
	}
	// Auto-renew TTL
	item.ExpiresAt = time.Now().Add(c.defaultTTL)
	c.cache.Add(key, item)
	return item.Value, true
}

func (c *ChronoCache) Delete(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cache.Remove(key)
}

func (c *ChronoCache) Snapshot() []SnapshotItem {
	c.lock.Lock()
	defer c.lock.Unlock()
	var snapshot []SnapshotItem
	now := time.Now()
	for _, k := range c.cache.Keys() {
		if item, ok := c.cache.Peek(k); ok {
			ttl := item.ExpiresAt.Sub(now)
			snapshot = append(snapshot, SnapshotItem{
				Key:     k,
				Value:   item.Value,
				TTLLeft: ttl,
			})
		}
	}
	return snapshot
}
