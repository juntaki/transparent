// Package transparent implements transparent cache operation.
package transparent

import "time"

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
	cache  BackendCache
	next   *Cache
	log    chan keyValue
	sync   chan bool
	synced chan bool
	done   chan bool
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
	c.sync = make(chan bool, 1)
	c.synced = make(chan bool, 1)

	go c.flush()
}

func (c *Cache) flush() {
	log := make(map[interface{}]interface{})
	done := false
	for { // infinite loop
	dedup:
		for { // loop for dedup request
			select {
			case kv, ok := <-c.log:
				// dedup the same key request
				if !ok {
					// channel is closed by Finalize
					done = true
					break dedup
				}
				log[kv.key] = kv.value

				// Too much keys cached
				if len(log) > 5 {
					break dedup
				}
			case <-c.sync:
				// Flush current buffer
				for k, v := range log {
					c.next.SetWriteBack(k, v)
				}
				// Flush value in channel
				// Switch to new channel for current writer
				old := *c
				c.log = make(chan keyValue, len(c.log))

				// Finalize old log
				close(old.log)
				old.flush()

				// Next, recursively
				if old.next != nil {
					old.next.sync <- true
					<-old.next.synced
				}

				c.synced <- true
			case <-time.After(time.Second * 1):
				// Flush if silent for one sec
				break dedup
			}
		}
		// flush value
		for k, v := range log {
			c.next.SetWriteBack(k, v)
		}
		if done {
			c.done <- true
			return
		}
		// reset
		log = make(map[interface{}]interface{})
	}
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
		c.log <- keyValue{key, value}
		c.next.SetWriteThrough(key, value)
	} else {
		c.log <- keyValue{key, value}
	}

	return
}

// Sync current buffered value
func (c *Cache) Sync() {
	c.sync <- true
	<-c.synced
}
