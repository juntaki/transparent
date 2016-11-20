// Package transparent implements transparent cache operation.
package transparent

// BackendCache supposes to be on-memory cache like LRU, or database, etc..
type BackendCache interface {
	Get(key interface{}) (interface{}, bool)
	Add(key interface{}, value interface{}) bool // Add key-value to cache
}

// Consider the following case
// [Backend cache] -> [Next cache] -> [Source]
//                                    ^
// [Another cache] ------------------/

// Cache is transparent interface to its backend cache
// Cache itself have CacheOps interface
type Cache struct {
	cache BackendCache
	next  *Cache
	log   chan keyValue
	done  chan bool
}

// Async log writer use this struct in its channel
type keyValue struct {
	key   interface{}
	value interface{}
}

// Initialize start goroutine for asynchronously set value
func (c *Cache) Initialize(size int) {
	c.log = make(chan keyValue, size)
	c.done = make(chan bool, 1)
	go func(c *Cache) {
		for {
			if kv, ok := <-c.log; ok {
				c.next.SetWriteBack(kv.key, kv.value)
			} else {
				c.done <- true
				return
			}
		}
	}(c)
}

// Finalize stops goroutine
func (c *Cache) Finalize() {
	close(c.log)
	<-c.done
}

// Get value from cache, or if not found, from source.
func (c *Cache) Get(key interface{}) interface{} {
	// Try to get backend cache
	value, found := c.cache.Get(key)
	if !found {
		// Recursively get value from source.
		value := c.next.Get(key)
		c.SetWriteBack(key, value)
		return value
	}
	return value
}

// SetWriteBack new value to Backend cache.
func (c *Cache) SetWriteBack(key interface{}, value interface{}) {
	c.setValue(key, value, false)
}

// SetWriteThrough set the value to Backend cache, Next cache, and Source
func (c *Cache) SetWriteThrough(key interface{}, value interface{}) {
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
		c.next.SetWriteThrough(key, value)
	} else {
		c.log <- keyValue{key, value}
	}

	return
}
