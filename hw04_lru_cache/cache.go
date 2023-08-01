package hw04lrucache

type Key string

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
	if exist {
		item.Value = value
		c.queue.MoveToFront(item)
	} else {
		if c.queue.Len() >= c.capacity {
			c.queue.Remove(c.queue.Back())
		}
		c.items[key] = c.queue.PushFront(value)
	}
	return exist
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	item, exist := c.exist(key)
	if exist {
		return item.Value, exist
	}
	return nil, exist
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
