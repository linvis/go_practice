package gocache

import (
	"errors"
	"sync"
)

var CacheErrorNotExist = errors.New("cache key not exist")
var CacheErrorEmpty = errors.New("cache empty")

type Cache struct {
	mutex   sync.Mutex
	maxByte int64
	cache   *LRUCache
}

func NewCache(maxByte int64) *Cache {
	return &Cache{
		maxByte: maxByte,
	}
}

func (c *Cache) Add(key string, val interface{}) {
	c.mutex.Lock()

	if c.cache == nil {
		c.cache = NewLRUCache(c.maxByte)
	}

	c.cache.Add(key, val)

	c.mutex.Unlock()
}

func (c *Cache) Get(key string) (interface{}, error) {
	if c.cache == nil {
		return nil, CacheErrorEmpty
	}
	return c.cache.Get(key)
}
