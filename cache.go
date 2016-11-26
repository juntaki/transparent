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
	*stacker
	BackendCache BackendCache // Target cache
	log          chan log     // Channel buffer
	sync         chan bool    // Control for flush buffer
	synced       chan bool
	done         chan bool
}

type message int

const (
	set message = iota
	remove
)

// Flush buffer use this struct in its log channel
type log struct {
	key interface{}
	*operation
}

type operation struct {
	value   interface{}
	message message
}

// New returns Cache layer.
func New(bufferSize int) *Cache {
	c := &Cache{
		log:    make(chan log, bufferSize),
		done:   make(chan bool, 1),
		sync:   make(chan bool, 1),
		synced: make(chan bool, 1),
	}
	c.stacker = &stacker{this: c}
	return c
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
	buffer := make(map[interface{}]operation)
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
				buffer[kv.key] = operation{kv.value, kv.message}

				// Too much keys cached
				if len(buffer) > 5 {
					break dedup
				}
			case <-c.sync:
				// Flush current buffer
				for k, v := range buffer {
					switch v.message {
					case remove:
						c.lower.Remove(k)
					case set:
						c.lower.Set(k, v.value)
					}
				}
				buffer = make(map[interface{}]operation)

				// Flush value in channel buffer
				//  Switch to new channel for current writer
				old := *c
				c.log = make(chan log, len(c.log))

				//  partially finalize old log for flushing
				close(old.log)
				old.flusher()
				<-old.done

				// Lower, recursively
				if old.lower != nil {
					old.lower.Sync()
				}

				c.synced <- true
			case <-time.After(time.Second * 1):
				// Flush if silent for one sec
				break dedup
			}
		}
		// Flush bufferd value
		for k, v := range buffer {
			switch v.message {
			case remove:
				c.lower.Remove(k)
			case set:
				c.lower.Set(k, v.value)
			}
		}
		// StopFlusher
		if done {
			c.done <- true
			return
		}
		// Reset buffer
		buffer = make(map[interface{}]operation)
	}
}

// Get value from cache, or if not found, recursively get.
func (c *Cache) Get(key interface{}) (value interface{}) {
	// Try to get backend cache
	value, found := c.BackendCache.Get(key)
	if !found {
		// Recursively get value from list.
		value := c.lower.Get(key)
		c.Set(key, value)
		return value
	}
	return value
}

// Set set new value to BackendCache.
func (c *Cache) Set(key interface{}, value interface{}) {
	if c.upper != nil {
		c.SkimOff(key)
	}
	c.BackendCache.Add(key, value)
	if c.lower == nil {
		// This backend cache is final destination
		return
	}
	// Queue to flush
	c.log <- log{key, &operation{value: value, message: set}}
	return
}

// Sync current buffered value
func (c *Cache) Sync() {
	c.sync <- true
	<-c.synced
}

// SkimOff remove upper layer's old value
func (c *Cache) SkimOff(key interface{}) {
	c.BackendCache.Remove(key)
	if c.upper == nil {
		// This is top layer
		return
	}
	c.upper.SkimOff(key)
}

// Remove recursively remove lower layer's value
func (c *Cache) Remove(key interface{}) {
	c.BackendCache.Remove(key)
	if c.lower == nil {
		// This is bottom layer
		return
	}
	// Queue to flush
	c.log <- log{key, &operation{nil, remove}}
	return
}
