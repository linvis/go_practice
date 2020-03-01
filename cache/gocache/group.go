package gocache

import "sync"

type GroupCallback func(key string) (interface{}, error)

type Group struct {
	Name     string
	callback GroupCallback
	cache    *Cache
}

type groupPool struct {
	groups map[string]*Group
	mutex  sync.RWMutex
}

var pool = groupPool{
	groups: make(map[string]*Group),
}

func NewGroup(name string, size int64, callback GroupCallback) *Group {
	g := &Group{
		Name:     name,
		callback: callback,
		cache:    NewCache(size),
	}

	pool.mutex.Lock()
	pool.groups[name] = g
	pool.mutex.Unlock()

	return g
}

func (g *Group) Get(key string) (interface{}, error) {
	v, err := g.cache.Get(key)
	if err != nil {
		return g.load(key)
	}

	return v, nil
}

func (g *Group) Set(key string, val interface{}) {
	g.cache.Add(key, val)
}

func (g *Group) load(key string) (interface{}, error) {
	v, err := g.callback(key)
	if err != nil {
		return nil, err
	}

	g.cache.Add(key, v)

	return v, nil
}

func GetGroup(name string) *Group {
	pool.mutex.RLock()

	defer pool.mutex.RUnlock()

	g, ok := pool.groups[name]
	if !ok {
		return nil
	}

	return g
}
