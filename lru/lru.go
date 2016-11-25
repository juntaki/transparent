// Package lru is simple and fast LRU implementation, for single thread.
// Cache is compatible LRU for transparent.BackendCache
package lru

// Cache is compatible LRU for transparent.BackendCache
type Cache struct {
	hash           map[interface{}]*keyValue
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
	}

	c.listHead.next = c.listHead
	c.listHead.prev = c.listHead
	return c
}

// Get value from cache if exist
func (c *Cache) Get(key interface{}) (value interface{}, found bool) {
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
func (c *Cache) Add(key interface{}, value interface{}) {
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

// Remove
func (c *Cache) Remove(key interface{}) {
	if kv, ok := c.hash[key]; ok {
		delete(c.hash, kv)
		listRemove(kv)
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
