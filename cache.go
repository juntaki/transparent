package transparent

import (
	"errors"
	"time"
)

// Cache provides operation of TransparentCache
type Cache struct {
	Storage Storage   // Target cache
	log     chan log  // Channel buffer
	sync    chan bool // Control for flush buffer
	synced  chan bool
	done    chan bool
	upper   Layer
	lower   Layer
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

// NewCache returns Cache layer.
func NewCache(bufferSize int, storage Storage) (*Cache, error) {
	if storage == nil {
		return nil, errors.New("empty storage")
	}
	c := &Cache{
		log:     make(chan log, bufferSize),
		done:    make(chan bool, 1),
		sync:    make(chan bool, 1),
		synced:  make(chan bool, 1),
		Storage: storage,
	}
	return c, nil
}

// DeleteCache clean up
func DeleteCache(c *Cache) {
	c.stopFlusher()
}

// StartFlusher starts flusher
func (c *Cache) startFlusher() {
	go c.flusher()
}

// StopFlusher stops flusher
func (c *Cache) stopFlusher() {
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
func (c *Cache) Get(key interface{}) (value interface{}, err error) {
	// Try to get backend cache
	value, err = c.Storage.Get(key)
	if err != nil {
		if c.lower == nil {
			return nil, errors.New("value not found")
		}
		// Recursively get value from list.
		value, err = c.lower.Get(key)
		if err != nil {
			return nil, err
		}
		err = c.Storage.Add(key, value)
		if err != nil {
			return nil, err
		}
	}
	return value, nil
}

// Set set new value to Storage.
func (c *Cache) Set(key interface{}, value interface{}) (err error) {
	err = c.Storage.Add(key, value)
	if err != nil {
		return err
	}
	if c.lower == nil {
		// This backend cache is final destination
		return nil
	}
	// Queue to flush
	c.log <- log{key, &operation{value: value, message: set}}
	return nil
}

// Sync current buffered value
func (c *Cache) Sync() error {
	c.sync <- true
	<-c.synced
	return nil
}

// Remove recursively remove lower layer's value
func (c *Cache) Remove(key interface{}) (err error) {
	err = c.Storage.Remove(key)
	if err != nil {
		return err
	}
	if c.lower == nil {
		// This is bottom layer
		return nil
	}
	// Queue to flush
	c.log <- log{key, &operation{nil, remove}}
	return nil
}

// SetUpper set upper layer
func (c *Cache) setUpper(upper Layer) {
	c.upper = upper
}

// SetLower set lower layer
func (c *Cache) setLower(lower Layer) {
	c.lower = lower
}
