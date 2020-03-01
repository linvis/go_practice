package gocache

import (
	"container/list"
	"errors"
)

var LRUInvalidKey = errors.New("Invalid Key")

type Cache struct {
	cache   map[string]*list.Element
	maxSize int
	list    *list.List
}

func New(size int) *Cache {
	return &Cache{
		cache:   make(map[string]*list.Element),
		maxSize: size,
		list:    list.New(),
	}
}

func (c *Cache) Add(key string, val interface{}) {
	ele, ok := c.cache[key]
	if ok {
		ele.Value = val
		c.list.MoveToBack(ele)
		return
	}

	if c.list.Len() >= c.maxSize {
		head := c.list.Front()
		c.list.Remove(head)
	}

	newEle := c.list.PushBack(val)
	c.cache[key] = newEle
}

func (c *Cache) Front() (interface{}, error) {
	if c.Len() <= 0 {
		return nil, errors.New("empty cache")
	}

	return c.list.Front().Value, nil
}

func (c *Cache) Back() (interface{}, error) {
	if c.Len() <= 0 {
		return nil, errors.New("empty cache")
	}

	return c.list.Back().Value, nil
}

func (c *Cache) Get(key string) (interface{}, error) {
	ele, ok := c.cache[key]
	if !ok {
		return nil, LRUInvalidKey
	}

	c.list.MoveToBack(ele)

	return ele.Value, nil
}

func (c *Cache) Delete(key string) error {
	ele, ok := c.cache[key]
	if !ok {
		return LRUInvalidKey
	}

	c.list.Remove(ele)
	delete(c.cache, key)

	return nil
}

func (c *Cache) Len() int {
	return c.list.Len()
}
