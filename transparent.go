// Package transparent A transparent package is a library that provides
// transparent caching operations for key-value stores. As shown in the figure
// below, it is possible to use relatively fast cache like LRU and slow
// and reliable storage like S3 via TransparentCache.
// Transparent Cache is tearable. In addition to caching, it is also possible to
// transparently use a layer of synchronization between distributed systems.
// See subpackage for example.
//
//  [Application]
//    |
//    v Get/Set
//  [Transparent cache] -[Flush buffer]-> [Next cache]
//   `-[Backend cache]                     `-[Source cache]
//      `-[LRU]                               `-[S3]
package transparent

import (
	"time"
)

// BackendCache defines the interface that TC's backend data storage destination should have.
// Both Get and Add should not be failed.
type BackendCache interface {
	Get(key interface{}) (value interface{}, found bool)
	Add(key interface{}, value interface{})
}

// Cache provides operation of TC
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

// New returns Cache, you can set nil to next, if it's Source.
func New(cache BackendCache, next *Cache, bufferSize int) *Cache {
	if cache == nil {
		return nil
	}
	return &Cache{
		cache:  cache,
		next:   next,
		log:    make(chan keyValue, bufferSize),
		done:   make(chan bool, 1),
		sync:   make(chan bool, 1),
		synced: make(chan bool, 1),
	}
}

// Initialize starts flusher
func (c *Cache) Initialize() {
	go c.flush()
}

// Finalize stops flusher
func (c *Cache) Finalize() {
	close(c.log)
	<-c.done
}

// flusher
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

// Get value from cache, or if not found, recursively get.
func (c *Cache) Get(key interface{}) (value interface{}) {
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

// SetWriteBack set new value to BackendCache.
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

// SetWriteThrough set the value to BackendCache and sync Source.
func (c *Cache) SetWriteThrough(key interface{}, value interface{}) {
	c.SetWriteBack(key, value)
	c.Sync()
}

// Sync current buffered value
func (c *Cache) Sync() {
	c.sync <- true
	<-c.synced
}
