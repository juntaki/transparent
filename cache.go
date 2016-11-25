package transparent

import "time"

// BackendCache defines the interface that TransparentCache's
// backend data storage destination should have.
// Add should not be failed.
type BackendCache interface {
	Get(key interface{}) (value interface{}, found bool)
	Add(key interface{}, value interface{})
	Remove(key interface{})
}

// Cache provides operation of TransparentCache
type Cache struct {
	BackendCache BackendCache // Target cache
	listHead     listHead
	log          chan keyValue // Channel buffer
	sync         chan bool     // Control for flush buffer
	synced       chan bool
	done         chan bool
}

// Flush buffer use this struct in its log channel
type keyValue struct {
	key   interface{}
	value interface{}
}

// New returns Cache layer.
func New(bufferSize int) *Cache {
	return &Cache{
		log:    make(chan keyValue, bufferSize),
		done:   make(chan bool, 1),
		sync:   make(chan bool, 1),
		synced: make(chan bool, 1),
	}
}

// StartFlusher starts flusher
func (c *Cache) StartFlusher() {
	go c.flusher()
}

// StopFlusher stops flusher
func (c *Cache) StopFlusher() {
	close(c.log)
	<-c.done
}

// Flusher
func (c *Cache) flusher() {
	buffer := make(map[interface{}]interface{})
	done := false
	for { // main loop
	dedup:
		for { // dedup request
			select {
			case kv, ok := <-c.log:
				if !ok {
					// channel is closed by StopFlusher
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
					c.listHead.lower.Set(k, v)
				}
				buffer = make(map[interface{}]interface{})

				// Flush value in channel buffer
				//  Switch to new channel for current writer
				old := *c
				c.log = make(chan keyValue, len(c.log))

				//  partially finalize old log for flushing
				close(old.log)
				old.flusher()
				<-old.done

				// Lower, recursively
				if old.listHead.lower != nil {
					old.listHead.lower.Sync()
				}

				c.synced <- true
			case <-time.After(time.Second * 1):
				// Flush if silent for one sec
				break dedup
			}
		}
		// Flush bufferd value
		for k, v := range buffer {
			c.listHead.lower.Set(k, v)
		}
		// StopFlusher
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
	value, found := c.BackendCache.Get(key)
	if !found {
		// Recursively get value from list.
		value := c.listHead.lower.Get(key)
		c.Set(key, value)
		return value
	}
	return value
}

// Set set new value to BackendCache.
func (c *Cache) Set(key interface{}, value interface{}) {
	c.BackendCache.Add(key, value)
	if c.listHead.lower == nil {
		// This backend cache is final destination
		return
	}
	// Queue to flush
	c.log <- keyValue{key, value}

	return
}

// Sync current buffered value
func (c *Cache) Sync() {
	c.sync <- true
	<-c.synced
}

// Remove
func (c *Cache) Remove(key interface{}) {
	c.BackendCache.Remove(key)
}

// Stack
func (c *Cache) Stack(l Layer) {
	c.getListHead().upper = l
	l.getListHead().lower = c
}

func (c *Cache) getListHead() *listHead {
	return &c.listHead
}
