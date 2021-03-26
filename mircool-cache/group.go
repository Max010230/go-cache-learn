package mircool_cache

import (
	"fmt"
	"go-cache-learn/mircool-cache/singleflight"
	"log"
	"sync"
)

type Group struct {
	name      string
	getter    Getter
	mainCache cache
	peers     PeerPicker
	loader    *singleflight.Group
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPick called more than once")
	}
	g.peers = peers
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (CacheData, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return CacheData{}, err
	}
	return CacheData{b: bytes}, err
}

func NewGroup(name string, memory int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	group := &Group{
		name:   name,
		getter: getter,
		mainCache: cache{
			memory: memory,
		},
		loader: &singleflight.Group{},
	}
	groups[name] = group
	return group
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	group := groups[name]
	return group
}

func (g *Group) Get(key string) (CacheData, error) {
	if key == "" {
		return CacheData{}, fmt.Errorf("key is required")
	}
	if data, ok := g.mainCache.Get(key); ok {
		log.Println("[cache]hit")
		return data, nil
	}
	return g.load(key)
}

func (g *Group) populateCache(key string, value CacheData) {
	g.mainCache.Add(key, value)
}

func (g *Group) getLocally(key string) (CacheData, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return CacheData{}, err
	}
	data := CacheData{
		b: cloneBytes(bytes),
	}
	g.populateCache(key, data)
	return data, err
}

func (g *Group) load(key string) (value CacheData, err error) {
	view, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err := g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[mircool cache] Failed to get from peer", err)
			}
		}
		return g.getLocally(key)
	})
	if err == nil {
		return view.(CacheData), nil
	}
	return

}
