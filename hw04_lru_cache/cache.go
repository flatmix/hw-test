package hw04lrucache

type Key string

type ValueStruct struct {
	Value interface{}
	Key   Key
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c lruCache) exist(key Key) (*ListItem, bool) {
	item, ok := c.items[key]
	return item, ok
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	item, exist := c.exist(key)
	NewElem := ValueStruct{
		Value: value,
		Key:   key,
	}
	if exist {
		item.Value = NewElem
		c.queue.MoveToFront(item)
	} else {
		if c.queue.Len() >= c.capacity {
			back := c.queue.Back()
			backElem, ok := back.Value.(ValueStruct)
			if ok {
				c.queue.Remove(back)
				delete(c.items, backElem.Key)
			}
		}
		c.items[key] = c.queue.PushFront(NewElem)
	}
	return exist
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	item, exist := c.exist(key)
	if exist {
		elem, ok := item.Value.(ValueStruct)
		if ok {
			if c.queue.Len() > 1 {
				c.queue.MoveToFront(item)
			}
			return elem.Value, exist
		}
	}
	return nil, exist
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
