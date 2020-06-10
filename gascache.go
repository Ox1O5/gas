package gas

import "sync"

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
	
}