package cache

import (
	"container/list"
	"fmt"
)

type Cache struct {
	storage  map[string]CacheEntry
	capacity int
	keyList  *list.List
}

type CacheEntry struct {
	key   string
	value string
}

func NewCache(capacity int) Cache {
	return Cache{
		storage:  make(map[string]CacheEntry),
		capacity: capacity,
		keyList:  list.New(),
	}
}

func (c *Cache) Get(key string) (string, error) {
	entry, ok := c.storage[key]
	if ok {
		element := c.keyList.Front()
		for element.Value != key {
			element = element.Next()
		}
		c.keyList.MoveToBack(element)
		return entry.value, nil
	}

	return "", fmt.Errorf("key %s not found", key)
}

func (c *Cache) Set(key, value string) {
	c.storage[key] = CacheEntry{
		key,
		value,
	}

	c.keyList.PushBack(key)

	if c.keyList.Len() > c.capacity {
		if evicted, ok := c.keyList.Remove(c.keyList.Front()).(string); ok {
			delete(c.storage, evicted)
		}
	}
}

func (c *Cache) Keys() []string {
	keys := []string{}

	for k, _ := range c.storage {
		keys = append(keys, k)
	}

	return keys
}
