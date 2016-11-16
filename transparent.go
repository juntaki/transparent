// Package transparent implements transparent cache operation.
package transparent

// BackendCache supposes to be on-memory cache like LRU, or database, etc..
type BackendCache interface {
	Get(key interface{}) (interface{}, bool)
	Add(key interface{}, value interface{}) bool // Add key-value to cache
}

// Consider the follwoing case
// [Backend cache] -> [Next cache] -> [Source]
//                                    ^
// [Another cache] ------------------/

// Cache is transparent interface to its backend cache
// Cache itself have CacheOps interface
type Cache struct {
	cache BackendCache
	next  *Cache
}

// Get value from cache, or if not found, from source.
func (c *Cache) Get(key interface{}) interface{} {
	// Try to get backend cache
	value, found := c.cache.Get(key)
	if !found {
		// Recursively get value from source.
		value := c.next.Get(key)
		c.Set(key, value)
		return value
	}
	return value
}

// Set new value to Backend cache.
func (c *Cache) Set(key interface{}, value interface{}) {
	c.setValue(key, value, false)
}

// SetSync set the value to Backend cache, Next cache, and Source
func (c *Cache) SetSync(key interface{}, value interface{}) {
	c.setValue(key, value, true)
}

func (c *Cache) setValue(key interface{}, value interface{}, sync bool) {
	c.cache.Add(key, value)

	if c.next == nil {
		// This backend is final destination
		return
	}

	// set value recursively
	if sync {
		c.next.SetSync(key, value)
	} else {
		go c.next.SetSync(key, value)
	}

	return
}

// SetWorld means SetSource + ensure Anoter cache is also up to date
func (c *Cache) SetWorld(key interface{}, value interface{}) bool {
	//TODO
	return false
}
