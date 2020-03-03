package gocache

import (
	"sync"

	pb "gocache/gocachepb"

	"github.com/golang/groupcache/singleflight"
)

type GroupCallback func(key string) ([]byte, error)

type Group struct {
	Name     string
	callback GroupCallback
	cache    *Cache
	peers    PeerPicker
	loader   *singleflight.Group
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
		loader:   &singleflight.Group{},
	}

	pool.mutex.Lock()
	pool.groups[name] = g
	pool.mutex.Unlock()

	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	v, err := g.cache.Get(key)
	if err != nil {
		return g.load(key)
	}

	return v, nil
}

func (g *Group) Set(key string, val ByteView) {
	g.cache.Add(key, val)
}

func (g *Group) loadFromLocal(key string) (ByteView, error) {
	v, err := g.callback(key)
	if err != nil {
		return ByteView{}, err
	}

	g.cache.Add(key, ByteView{b: v})

	return ByteView{b: v}, nil
}

func (g *Group) load(key string) (ByteView, error) {
	// if g.peers != nil {
	// 	peer, err := g.peers.PeerPick(key)
	// 	if err == nil {
	// 		return g.GetFromPeer(peer, key)
	// 	}
	// }

	// return g.loadFromLocal(key)

	view, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			peer, err := g.peers.PeerPick(key)
			if err == nil {
				return g.GetFromPeer(peer, key)
			}
		}

		return g.loadFromLocal(key)
	})

	if err == nil {
		return view.(ByteView), nil
	} else {
		return ByteView{}, nil
	}
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

func (g *Group) RegisterPeer(peers PeerPicker) {
	if g.peers != nil {
		panic("only one peer is allowed")
	}

	g.peers = peers
}

func (g *Group) GetFromPeer(peer PeerGetter, key string) (ByteView, error) {
	req := &pb.Request{
		Group: g.Name,
		Key:   key,
	}
	res := &pb.Response{}

	err := peer.Get(req, res)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: res.Value}, nil
}
