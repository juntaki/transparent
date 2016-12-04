// Package lru is simple LRU implementation.
// Cache is compatible for transparent.Storage
package lru

import (
	"errors"
	"sync"
)

// Cache is compatible LRU for transparent.BackendCache
type Cache struct {
	hash           map[interface{}]*keyValue
	lock           sync.RWMutex
	listHead       *keyValue
	currentEntries int
	maxEntries     int
}

type keyValue struct {
	key   interface{}
	value interface{}
	prev  *keyValue
	next  *keyValue
}

// New returns empty LRU Cache
func New(maxEntries int) *Cache {
	c := &Cache{
		hash:           make(map[interface{}]*keyValue),
		currentEntries: 0,
		maxEntries:     maxEntries,
		listHead:       &keyValue{},
		lock:           sync.RWMutex{},
	}

	c.listHead.next = c.listHead
	c.listHead.prev = c.listHead
	return c
}

// Get value from cache if exist
func (c *Cache) Get(key interface{}) (value interface{}, err error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if kv, ok := c.hash[key]; ok {
		if kv != c.listHead.next {
			listRemove(kv)
			listAdd(c.listHead, kv)
		}
		return kv.value, nil
	}
	return nil, errors.New("value not found")
}

// Add value to cache
func (c *Cache) Add(key interface{}, value interface{}) (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if kv, ok := c.hash[key]; ok {
		if kv != c.listHead.next {
			listRemove(kv)
			listAdd(c.listHead, kv)
		}
		kv.value = value
	} else {
		if c.maxEntries != c.currentEntries {
			c.currentEntries++
		} else {
			lastItem := c.listHead.prev
			delete(c.hash, lastItem.key)
			listRemove(lastItem)
		}

		kv := &keyValue{
			key:   key,
			value: value,
		}
		listAdd(c.listHead, kv)
		c.hash[key] = kv
	}
	return nil
}

// Remove value from cache
func (c *Cache) Remove(key interface{}) (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if kv, ok := c.hash[key]; ok {
		delete(c.hash, key)
		listRemove(kv)
	}
	return nil
}

func listRemove(kv *keyValue) {
	kv.prev.next = kv.next
	kv.next.prev = kv.prev
}
func listAdd(prev, kv *keyValue) {
	next := prev.next
	kv.next = next
	kv.prev = prev
	next.prev = kv
	prev.next = kv
}
