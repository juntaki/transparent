// Package transparent implements transparent cache operation.
package transparent

// BackendCache supposes to be on-memory cache like LRU, or database, etc..
type BackendCache interface {
	Get(key interface{}) (interface{}, bool)
	Add(key interface{}, value interface{}) bool // Add key-value to cache
}

// CacheOps supports Get and multiple type of Set
type CacheOps interface {
	// Get value from cache, or if not found, from source.
	Get(key interface{}) (interface{}, bool)

	// Consider the follwoing case
	// [Backend cache] -> [Next cache] -> [Source]
	//                                    ^
	// [Another cache] ------------------/

	// Set new value to Backend cache only
	Set(key interface{}, value interface{}) bool

	// SetSource set the value to Backend cache, Next cache, and Source
	SetSource(key interface{}, value interface{}) bool

	// SetSource + ensure Anoter cache is also up to date
	SetWorld(key interface{}, value interface{}) bool
}

// Cache is transparent interface to its backend cache
// Cache itself have CacheOps interface
type Cache struct {
	cache BackendCache
	next  CacheOps
}

// Get value from cache, or if not found, from source.
func (c Cache) Get(key interface{}) (interface{}, bool) {
	// Try to get backend cache
	value, ok := c.cache.Get(key)
	if !ok {
		// Recursively get value from source.
		value, ok := c.next.Get(key)
		if !ok {
			return nil, false
		}
		c.Set(key, value)
		return value, true
	}
	return value, true
}

// Set new value to Backend cache only
func (c Cache) Set(key interface{}, value interface{}) bool {
	return c.cache.Add(key, value)
}

// SetSource set the value to Backend cache, Next cache, and Source
func (c Cache) SetSource(key interface{}, value interface{}) bool {
	ok := c.Set(key, value)
	if !ok {
		return false
	}

	if c.next == nil {
		// This backend is final destination
		return true
	}

	// set value recursively
	c.next.SetSource(key, value)
	return true
}

// SetWorld means SetSource + ensure Anoter cache is also up to date
func (c Cache) SetWorld(key interface{}, value interface{}) bool {
	//TODO
	return false
}
