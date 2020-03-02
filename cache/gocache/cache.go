package gocache

import (
	"errors"
	"gocache/lru"
	"sync"
)

var CacheErrorNotExist = errors.New("cache key not exist")
var CacheErrorEmpty = errors.New("cache empty")

type Cache struct {
	mutex   sync.Mutex
	maxByte int64
	cache   *lru.Cache
}

func NewCache(maxByte int64) *Cache {
	return &Cache{
		maxByte: maxByte,
	}
}

func (c *Cache) Add(key string, val ByteView) {
	c.mutex.Lock()

	if c.cache == nil {
		c.cache = lru.New(c.maxByte)
	}

	c.cache.Add(key, val)

	c.mutex.Unlock()
}

func (c *Cache) Get(key string) (ByteView, error) {
	if c.cache == nil {
		return ByteView{}, CacheErrorEmpty
	}
	v, err := c.cache.Get(key)

	return v.(ByteView), err
}
