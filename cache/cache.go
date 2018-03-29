package cache

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	storage        map[string]CacheEntry
	capacity       int
	expirationTime time.Duration
	keyList        *list.List
	mutex          *sync.Mutex
}

type CacheEntry struct {
	key          string
	value        string
	creationTime time.Time
}

func NewCache(capacity int, expire time.Duration) *Cache {
	c := &Cache{
		storage:        make(map[string]CacheEntry),
		capacity:       capacity,
		expirationTime: expire,
		keyList:        list.New(),
		mutex:          &sync.Mutex{},
	}

	go c.Start()
	return c
}

func (c *Cache) Get(key string) (string, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, ok := c.storage[key]
	if ok {
		element := c.findElementForKey(key)
		c.keyList.MoveToBack(element)
		return entry.value, nil
	}

	return "", fmt.Errorf("key %s not found", key)
}

func (c *Cache) findElementForKey(key string) *list.Element {
	element := c.keyList.Front()
	for element.Value != key {
		element = element.Next()
	}

	return element
}

func (c *Cache) Set(key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.storage[key] = CacheEntry{
		key,
		value,
		time.Now(),
	}

	c.keyList.PushBack(key)

	if c.keyList.Len() > c.capacity {
		if evicted, ok := c.keyList.Remove(c.keyList.Front()).(string); ok {
			delete(c.storage, evicted)
		}
	}
}

func (c *Cache) Keys() []string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	keys := []string{}

	for k, _ := range c.storage {
		keys = append(keys, k)
	}

	return keys
}

func (c *Cache) Start() {
	for {
		c.mutex.Lock()

		keysToDelete := []string{}

		for k, v := range c.storage {
			if time.Since(v.creationTime) >= c.expirationTime {
				keysToDelete = append(keysToDelete, k)
			}
		}

		for _, key := range keysToDelete {
			e := c.findElementForKey(key)
			delete(c.storage, key)
			c.keyList.Remove(e)
		}

		time.Sleep(c.expirationTime / 10)
		c.mutex.Unlock()
	}
}
