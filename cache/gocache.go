package gocache

type GroupCallback func(key string) (interface{}, error)

type Group struct {
	Name     string
	callback GroupCallback
	cache    *Cache
}

func NewGroup(name string, size int64, callback GroupCallback) *Group {
	return &Group{
		Name:     name,
		callback: callback,
		cache:    NewCache(size),
	}
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
