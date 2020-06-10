package lru

import "container/list"

//并发不安全的最近最久未使用cache
type Cache struct {
	//maxBytes等于0表示没有限制
	maxBytes int64
	nbytes int64
	ll *list.List
	cache map[string]*list.Element
	// entry被移除的回调函数，（可选函数，可以为nil)
	OnEviced func(key string, value Value)
}

//设为接口方便扩展
type Value interface {
	Len() int
}

//一条记录,作为双链表list的结点，保存key方便从map里删除
type entry struct {
	key string
	value Value
}

//Cache构造函数
func New(maxBytes int64, onEviced func(string, Value)) *Cache  {
	return &Cache{
		maxBytes: maxBytes,
		ll: list.New(),
		cache: make(map[string]*list.Element),
		OnEviced: onEviced,
	}
}

//Add 增，新增记录，并移动到队尾
func (c *Cache)Add(key string, value Value) {
	if elem, hit := c.cache[key]; hit {
		c.ll.MoveToBack(elem)
		kv := elem.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		elem := c.ll.PushBack(&entry{key, value})
		c.cache[key] = elem
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	for c.maxBytes !=0 && c.nbytes > c.maxBytes {
		c.RemoveOldest()
	}
}

//Remove 删，移除最近最久未访问结点
func (c *Cache) RemoveOldest() {
	elem := c.ll.Front()
	if elem != nil {
		kv := elem.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		c.ll.Remove(elem)
		if c.OnEviced != nil {
			c.OnEviced(kv.key, kv.value)
		}
	}
}

//Get 查找链表中结点，并移动到队尾表示刚访问过
func (c *Cache)Get(key string) (value Value, ok bool) {
	if elem, hit := c.cache[key]; hit {
		c.ll.MoveToBack(elem)
		kv := elem.Value.(*entry)
		return kv.value, true
	}
	return
}

func (c *Cache)Len() int {
	return c.ll.Len()
}

