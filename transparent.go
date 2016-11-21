// Package transparent implements transparent cache operation.
package transparent

import "time"

// Key is comparable value
type Key interface{}

// BackendCache supposes to be on-memory cache like LRU, or database, etc..
type BackendCache interface {
	Get(key Key) (interface{}, bool)
	Add(key Key, value interface{})
}

//  [Application]
//    |
//    v Get/Set
//  [Transparent cache] -[Flush buffer]-> [Next cache]
//   `-[Backend cache]                     `-[Source cache]
//      `-[LRU]                               `-[S3]

// Cache is transparent interface to its backend cache
// You can stack Cache for tiring
type Cache struct {
	cache  BackendCache  // Target cache
	next   *Cache        // Next should be more stable but slow
	log    chan keyValue // Channel buffer
	sync   chan bool     // Control for flush buffer
	synced chan bool
	done   chan bool
}

// Flush buffer use this struct in its log channel
type keyValue struct {
	key   interface{}
	value interface{}
}

// Initialize start flush buffer goroutine for asynchronously set value
func (c *Cache) Initialize(bufferSize int) {
	c.log = make(chan keyValue, bufferSize)
	c.done = make(chan bool, 1)
	c.sync = make(chan bool, 1)
	c.synced = make(chan bool, 1)

	go c.flush()
}

// Flush buffer
func (c *Cache) flush() {
	buffer := make(map[interface{}]interface{})
	done := false
	for { // main loop
	dedup:
		for { // dedup request
			select {
			case kv, ok := <-c.log:
				if !ok {
					// channel is closed by Finalize
					done = true
					break dedup
				}
				buffer[kv.key] = kv.value

				// Too much keys cached
				if len(buffer) > 5 {
					break dedup
				}
			case <-c.sync:
				// Flush current buffer
				for k, v := range buffer {
					c.next.SetWriteBack(k, v)
				}
				buffer = make(map[interface{}]interface{})

				// Flush value in channel buffer
				//  Switch to new channel for current writer
				old := *c
				c.log = make(chan keyValue, len(c.log))

				//  partially finalize old log for flushing
				close(old.log)
				old.flush()
				<-old.done

				// Next, recursively
				if old.next != nil {
					old.next.Sync()
				}

				c.synced <- true
			case <-time.After(time.Second * 1):
				// Flush if silent for one sec
				break dedup
			}
		}
		// Flush bufferd value
		for k, v := range buffer {
			c.next.SetWriteBack(k, v)
		}
		// Finalize
		if done {
			c.done <- true
			return
		}
		// Reset buffer
		buffer = make(map[interface{}]interface{})
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
		// Recursively get value from list.
		value := c.next.Get(key)
		c.SetWriteBack(key, value)
		return value
	}
	return value
}

// SetWriteBack new value to Backend cache.
func (c *Cache) SetWriteBack(key interface{}, value interface{}) {
	c.cache.Add(key, value)
	if c.next == nil {
		// This backend cache is final destination
		return
	}
	// Queue to flush
	c.log <- keyValue{key, value}

	return
}

// SetWriteThrough set the value and sync
func (c *Cache) SetWriteThrough(key interface{}, value interface{}) {
	c.SetWriteBack(key, value)
	c.Sync()
}

// Sync current buffered value
func (c *Cache) Sync() {
	c.sync <- true
	<-c.synced
}
