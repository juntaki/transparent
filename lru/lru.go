// Package lru is simple LRU implementation.
// Cache is compatible for transparent.Storage
package lru

import (
	"sync"

	"github.com/juntaki/transparent"
)

type storage struct {
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

// NewStorage returns LRU Storage
func NewStorage(maxEntries int) transparent.BackendStorage {
	c := &storage{
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
func (c *storage) Get(key interface{}) (value interface{}, err error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if kv, ok := c.hash[key]; ok {
		if kv != c.listHead.next {
			listRemove(kv)
			listAdd(c.listHead, kv)
		}
		return kv.value, nil
	}
	return nil, &transparent.KeyNotFoundError{Key: key}
}

// Add value to cache
func (c *storage) Add(key interface{}, value interface{}) (err error) {
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
func (c *storage) Remove(key interface{}) (err error) {
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
