package mircool_cache

import (
	"go-cache-learn/mircool-cache/lru"
	"sync"
)

type cache struct {
	mu     sync.Mutex
	lru    *lru.Cache
	memory int64
}

func (c *cache) Add(key string, value CacheData) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.memory, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) Get(key string) (data CacheData, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if value, ok := c.lru.Get(key); ok {
		return value.(CacheData), true
	}
	return
}
