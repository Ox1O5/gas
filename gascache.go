package gas

import (
	"fmt"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc)Get(key string) ([]byte, error) {
	return f(key)
}


// A Group is a cache namespace and associated data loaded spread over
// a group of 1 or more machines.
type Group struct{
	name string
	getter 	Getter
	mainCache cache
}

var (
	mu sync.RWMutex
	groups = make(map[string]*Group)
)

//Constructor
func NewGroup(name string, cacheByte int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name: name,
		getter: getter,
		mainCache: cache{cacheBytes: cacheByte},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group)Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[Gascache] hit")
		return v, nil
	}

	return g.load(key)
}

func (g *Group)load(key string)(ByteView, error) {
	return g.getLocal(key)
}

func (g *Group)getLocal(key string)(ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b : byteClone(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func(g *Group)populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}