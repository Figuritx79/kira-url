package cache

import (
	"sync"

	"github.com/coocood/freecache"
)

type Cache struct {
	cache     *freecache.Cache
	hits      map[string]int
	hitsMutex sync.RWMutex
}

func NewCache(size int) *Cache {
	return &Cache{
		cache: freecache.NewCache(size),
		hits:  make(map[string]int),
	}
}

func (c *Cache) Set(key string, value []byte, baseTTL int) {
	c.hitsMutex.RLock()
	hits := c.hits[key]
	c.hitsMutex.RUnlock()

	// Extend TTL for frequently accessed keys
	ttl := baseTTL
	if hits > 10 {
		ttl *= 2 // Double TTL for popular items
	}
	c.cache.Set([]byte(key), value, ttl)
}

func (c *Cache) Get(key string) ([]byte, error) {
	val, err := c.cache.Get([]byte(key))
	if err == nil {
		c.hitsMutex.Lock()
		c.hits[key]++
		c.hitsMutex.Unlock()
	}
	return val, err
}
