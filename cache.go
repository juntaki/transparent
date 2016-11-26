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

type buffer struct {
	queue map[interface{}]operation
	c     *Cache
	limit int
}

func (b *buffer) reset() {
	b.queue = make(map[interface{}]operation)
}

func (b *buffer) add(l *log) {
	b.queue[l.key] = operation{l.value, l.message}
}

func (b *buffer) checkLimit() {
	if len(b.queue) > b.limit {
		b.flush()
	}
}

func (b *buffer) flush() {
	for k, o := range b.queue {
		switch o.message {
		case remove:
			b.c.lower.Remove(k)
		case set:
			b.c.lower.Set(k, o.value)
		}
	}
	b.reset()
}

// Flusher
func (c *Cache) flusher() {
	b := buffer{c: c, limit: 5}
	b.reset()
done:
	for { // main loop
		select {
		case l, ok := <-c.log:
			if !ok {
				// channel is closed by StopFlusher
				break done
			}
			b.add(&l)
			b.checkLimit()
		case <-c.sync:
			// Flush current buffer
			b.flush()

			// Flush value in channel buffer
			// Switch to new channel for current writer
			old := *c
			c.log = make(chan log, len(c.log))

			// Close old log for flushing
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
			b.flush()
		}
	}
	// Flush bufferd value
	b.flush()
	c.done <- true
	return
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
