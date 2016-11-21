// Package lru is simple and fast, for single thread.
package lru

import "github.com/juntaki/transparent"

// Cache is compatible LRU for transparent.BackendCache
type Cache struct {
	hash        map[transparent.Key]*keyValue
	listHead    *keyValue
	currentSize int
	limitSize   int
}

type keyValue struct {
	key   transparent.Key
	value interface{}
	prev  *keyValue
	next  *keyValue
}

// New returns empty LRU Cache
func New(size int) *Cache {
	c := &Cache{
		hash:        make(map[transparent.Key]*keyValue),
		currentSize: 0,
		limitSize:   size,
		listHead:    &keyValue{},
	}

	c.listHead.next = c.listHead
	c.listHead.prev = c.listHead
	return c
}

// Get value from cache if exist
func (c *Cache) Get(key transparent.Key) (interface{}, bool) {
	if kv, ok := c.hash[key]; ok {
		if kv != c.listHead.next {
			listRemove(kv)
			listAdd(c.listHead, kv)
		}
		return kv.value, true
	}
	return nil, false
}

// Add value to cache
func (c *Cache) Add(key transparent.Key, value interface{}) {
	if kv, ok := c.hash[key]; ok {
		if kv != c.listHead.next {
			listRemove(kv)
			listAdd(c.listHead, kv)
		}
		kv.value = value
	} else {
		if c.limitSize != c.currentSize {
			c.currentSize++
		} else {
			lastItem := c.listHead.prev
			delete(c.hash, lastItem)
			listRemove(lastItem)
		}

		kv := &keyValue{
			key:   key,
			value: value,
		}
		listAdd(c.listHead, kv)
		c.hash[key] = kv
	}
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
