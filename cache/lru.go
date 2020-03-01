package gocache

import (
	"container/list"
	"errors"
	"unsafe"
)

type LRUCache struct {
	cache   map[string]*list.Element
	maxByte int64
	nbytes  int64
	list    *list.List
}

type entry struct {
	key   string
	value interface{}
}

func NewLRUCache(maxByte int64) *LRUCache {
	return &LRUCache{
		cache:   make(map[string]*list.Element),
		maxByte: maxByte,
		list:    list.New(),
	}
}

func (c *LRUCache) Add(key string, val interface{}) {
	ele, ok := c.cache[key]
	if ok {
		entry := ele.Value.(*entry)
		entry.value = val
		c.list.MoveToBack(ele)
		return
	}

	needSize := int64(unsafe.Sizeof(key)) + int64(unsafe.Sizeof(val))
	if needSize > c.Free() {
		head := c.list.Front().Value.(*entry)

		c.Delete(head.key)

		c.Add(key, val)

		return
	}

	newEle := c.list.PushBack(&entry{key, val})
	c.cache[key] = newEle
}

func (c *LRUCache) Front() (interface{}, error) {
	if c.Len() <= 0 {
		return nil, errors.New("empty cache")
	}

	return c.list.Front().Value.(*entry).value, nil
}

func (c *LRUCache) Back() (interface{}, error) {
	if c.Len() <= 0 {
		return nil, errors.New("empty cache")
	}

	return c.list.Back().Value.(*entry).value, nil
}

func (c *LRUCache) Get(key string) (interface{}, error) {
	ele, ok := c.cache[key]
	if !ok {
		return nil, CacheErrorNotExist
	}

	c.list.MoveToBack(ele)

	return ele.Value.(*entry).value, nil
}

func (c *LRUCache) Delete(key string) error {
	ele, ok := c.cache[key]
	if !ok {
		return CacheErrorNotExist
	}

	c.list.Remove(ele)
	delete(c.cache, key)

	c.nbytes -= int64(unsafe.Sizeof(ele.Value))

	return nil
}

func (c *LRUCache) Len() int {
	return c.list.Len()
}

func (c *LRUCache) Size() int64 {
	return c.nbytes
}

func (c *LRUCache) Free() int64 {
	return c.maxByte - c.nbytes
}
