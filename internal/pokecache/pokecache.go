package pokecache

import (
	"time"
	"sync"
)

type Cache struct {
    cacheEntries map[string]CacheEntry
	interval time.Duration
	mux *sync.Mutex
}

type CacheEntry struct {
    createdAt time.Time
	value []byte
}

func NewCache(interval time.Duration) Cache {
	cache := Cache{
		cacheEntries: make(map[string]CacheEntry),
		interval: interval,
		mux: &sync.Mutex{},
	}

	go cache.reapLoop()
	return cache
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()
	entry, ok := c.cacheEntries[key]
	return entry.value, ok
}

func (c *Cache) Add(key string, value []byte) {
	newCacheEntry := CacheEntry{
		createdAt: time.Now(),
		value: value,
	}

	c.mux.Lock()
	defer c.mux.Unlock()
	c.cacheEntries[key] = newCacheEntry
}

func (c *Cache) reapLoop(){
	_ = time.NewTicker(c.interval)

	for {
		c.mux.Lock()
		for key, entry := range c.cacheEntries {
			elapsed := time.Since(entry.createdAt)
			if elapsed.Seconds() > c.interval.Seconds() {
				delete(c.cacheEntries, key)
			}
		}
		c.mux.Unlock()
	}

}
