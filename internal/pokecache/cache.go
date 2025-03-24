package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries  map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
}

func (c *Cache) Add(key string, val []byte) {
	entry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.entries[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	entry, ok := c.entries[key]
	if ok {
		return entry.val, ok
	}
	return []byte{}, false
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	for {
		<-ticker.C
		c.mu.Lock()
		threshold := time.Now().Add(-c.interval)
		for key, entry := range c.entries {
			if entry.createdAt.Before(threshold) {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		entries:  make(map[string]cacheEntry),
		interval: interval,
	}
	go cache.reapLoop()
	return cache
}
