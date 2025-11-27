package utils

import (
	"time"
	"sync"
)


type item[V any] struct {
	value V
	expiry time.Time
}

func (i item[V]) isExpired() bool {
	return time.Now().After(i.expiry)
}

type TTLCache[K comparable, V any] struct {
	items map[K]item[V]
	mu sync.Mutex
}

func NewTTLCache[K comparable, V any]() *TTLCache[K,V]{
	c := &TTLCache[K,V]{
		items: make(map[K]item[V]),
	}
	go func() {
		for range time.Tick(5 * time.Second) {
			c.mu.Lock()
			for key, item := range c.items {
				if item.isExpired() {
					delete(c.items, key)
				}
			}
			c.mu.Unlock()
		}
	}()
	return c
}

func (c *TTLCache[K,V]) Set(key K, value V, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = item[V]{
		value: value,
		expiry: time.Now().Add(ttl),
	}
}

func (c *TTLCache[K,V]) Get(key K) (V, bool){
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if !ok {
		return item.value, false
	}
	if item.isExpired() {
		delete(c.items, key)
		return item.value, false
	}
	return item.value, true
}
