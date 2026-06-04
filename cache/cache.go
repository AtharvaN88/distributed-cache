package cache

import (
	"container/list"
	"sync"
	"time"
)

type CacheItem struct {
	Value string
	Expiration int64
	Element *list.Element
}

type Cache struct {
	items map[string] CacheItem
	mu sync.RWMutex
	order *list.List
	capacity int
}

func NewCache(cleanupInterval int) *Cache {
	c := &Cache {
		items: make(map[string]CacheItem),
		order:    list.New(),
		capacity: capacity,
	}

	go func() {
		for {
			time.Sleep(time.Duration(cleanupInterval) * time.Second)
			c.cleanup()
		}
	}()

	return c
}

// Set a key with optional TTL
func (c *Cache) Set(key, value string, ttl int64) {
	c.mu.Lock()
	defer c.mu.Unlock()


	if item, ok := c.items[key]; ok {
		// Update value
		item.Value = value
		item.Expiry = time.Now().Add(time.Duration(ttl) * time.Second)
		// Move to front
		c.order.MoveToFront(item.Element)
		return
	}

	// Add new item
	item := &CacheItem{
		Value:  value,
		Expiry: time.Now().Add(time.Duration(ttl) * time.Second),
	}
	item.Element = c.order.PushFront(key)
	c.items[key] = item

	// Evict LRU if over capacity
	if c.order.Len() > c.capacity {
		c.evict()
	}
}

// Get a key, returns empty string if not found or expired
func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[key]
	if !ok || item.Expiry.Before(time.Now()) {
		return "", false
	}

	c.order.MoveToFront(item.Element)

	return item.Value, true
}

// cleanup expired items
func(c *Cache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()


	for k, v := range c.items {
		if v.Expiration > 0 && time.Now().Unix() > v.Expiration {
			delete(c.items, k)
		}
	}
}

func (c *Cache) evict() {
	// Remove last element
	back := c.order.Back()
	if back == nil {
		return
	}

	key := back.Value.(string)
	c.order.Remove(back)
	delete(c.items, key)
}