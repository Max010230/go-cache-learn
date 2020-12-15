package lru

import "container/list"

type Value interface {
	Len() int
}

type entry struct {
	key   string
	value Value
}

type Cache struct {
	maxBytes    int64                         //最大缓存空间
	usedBytes   int64                         //已使用缓存空间
	memorySpace *list.List                    //缓存存储的双向链表
	cacheSpace  map[string]*list.Element      //缓存存储的map
	OnEvicted   func(key string, value Value) //缓存删除回调
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:    maxBytes,
		memorySpace: list.New(),
		cacheSpace:  make(map[string]*list.Element),
		OnEvicted:   onEvicted,
	}
}

//缓存查询，获取到缓存之后，将缓存移动至队尾
func (c *Cache) Get(key string) (value Value, ok bool) {
	if element, ok := c.cacheSpace[key]; ok {
		c.memorySpace.MoveToBack(element)
		e := element.Value.(*entry)
		return e.value, true
	}
	return
}

//删除最近最少使用的缓存
//1、取到队首元素
//2、双向列表删除队首元素
//3、map删除对应的键值对
//4、重新计算已使用内存
//5、回调方法存在，调用回调方法
func (c *Cache) RemoveOldest() {
	front := c.memorySpace.Front()
	if front != nil {
		c.memorySpace.Remove(front)
		e := front.Value.(*entry)
		delete(c.cacheSpace, e.key)
		c.usedBytes -= int64(len(e.key)) + int64(e.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(e.key, e.value)
		}
	}
}

//添加缓存
func (c *Cache) Add(key string, value Value) {
	if element, ok := c.cacheSpace[key]; ok {
		c.memorySpace.MoveToBack(element)
		e := element.Value.(*entry)
		c.usedBytes += int64(value.Len()) - int64(e.value.Len())
		e.value = value
	} else {
		element := c.memorySpace.PushBack(&entry{
			key:   key,
			value: value,
		})
		c.cacheSpace[key] = element
		c.usedBytes += int64(len(key)) + int64(value.Len())
	}
	//防止超出存储上限
	for c.maxBytes < c.usedBytes && c.maxBytes != 0 {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.memorySpace.Len()
}
