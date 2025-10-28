package hw04lrucache

import "sync"

type Key string

// Хранит ключ и значение.
type cacheItem struct {
	key   Key
	value interface{}
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu       sync.RWMutex
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

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, exists := c.items[key]; exists {
		// Элемент уже существует. Обновляем значение и перемещаем в начало.
		item := elem.Value.(*cacheItem)
		item.value = value
		c.queue.MoveToFront(elem)
		return true
	}

	// Добавляем новый.
	newItem := &cacheItem{key: key, value: value}
	newListElem := c.queue.PushFront(newItem)
	c.items[key] = newListElem

	// Проверяем длину.
	if c.queue.Len() > c.capacity {
		back := c.queue.Back()
		if back != nil {
			backItem := back.Value.(*cacheItem)
			delete(c.items, backItem.key)
			c.queue.Remove(back)
		}
	}
	return false
}

// Получаем элемент по ключу.
func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	elem, exists := c.items[key]
	if !exists {
		return nil, false
	}

	// Перемещаем элемент в начало.
	c.queue.MoveToFront(elem)

	item := elem.Value.(*cacheItem)
	return item.value, true
}

// Очищаем кэш.
func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem)
}
