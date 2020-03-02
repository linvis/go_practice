package lru

import (
	"container/list"
	"errors"
	"unsafe"
)

type Cache struct {
	cache   map[string]*list.Element
	maxByte int64
	nbytes  int64
	list    *list.List
}

type Value interface {
	Len() int
}

type entry struct {
	key   string
	value Value
}

var CacheErrorNotExist = errors.New("cache key not exist")
var CacheErrorEmpty = errors.New("cache empty")

func New(maxByte int64) *Cache {
	return &Cache{
		cache:   make(map[string]*list.Element),
		maxByte: maxByte,
		list:    list.New(),
	}
}

func (c *Cache) Add(key string, val Value) {
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

func (c *Cache) Front() (Value, error) {
	if c.Len() <= 0 {
		return nil, errors.New("empty cache")
	}

	return c.list.Front().Value.(*entry).value, nil
}

func (c *Cache) Back() (Value, error) {
	if c.Len() <= 0 {
		return nil, errors.New("empty cache")
	}

	return c.list.Back().Value.(*entry).value, nil
}

func (c *Cache) Get(key string) (Value, error) {
	ele, ok := c.cache[key]
	if !ok {
		return nil, CacheErrorNotExist
	}

	c.list.MoveToBack(ele)

	return ele.Value.(*entry).value, nil
}

func (c *Cache) Delete(key string) error {
	ele, ok := c.cache[key]
	if !ok {
		return CacheErrorNotExist
	}

	c.list.Remove(ele)
	delete(c.cache, key)

	c.nbytes -= int64(unsafe.Sizeof(ele.Value))

	return nil
}

func (c *Cache) Len() int {
	return c.list.Len()
}

func (c *Cache) Size() int64 {
	return c.nbytes
}

func (c *Cache) Free() int64 {
	return c.maxByte - c.nbytes
}
