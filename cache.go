package transparent

import (
	"errors"
	"time"
)

// LayerCache wraps BackendStorage.
// It Get/Set key-value to BackendStorage,
// and asynchronously apply same operation to Next Layer.
// It must be Stacked on a Layer.
type LayerCache struct {
	Storage BackendStorage // Target cache
	log     chan log       // Channel buffer
	sync    chan bool      // Control for flush buffer
	synced  chan bool
	done    chan bool
	next    Layer
}

// Flush buffer use this struct in its log channel
type log struct {
	key interface{}
	*Message
}

// NewLayerCache returns Cache layer.
func NewLayerCache(bufferSize int, storage BackendStorage) (*LayerCache, error) {
	if storage == nil {
		return nil, errors.New("empty storage")
	}
	c := &LayerCache{
		log:     make(chan log, bufferSize),
		done:    make(chan bool, 1),
		sync:    make(chan bool, 1),
		synced:  make(chan bool, 1),
		Storage: storage,
	}
	return c, nil
}

func (c *LayerCache) start() error {
	go c.flusher()
	return nil
}

func (c *LayerCache) stop() error {
	close(c.log)
	<-c.done
	return nil
}

type buffer struct {
	queue map[interface{}]*Message
	c     *LayerCache
	limit int
}

func (b *buffer) reset() {
	b.queue = make(map[interface{}]*Message)
}

func (b *buffer) add(l *log) {
	b.queue[l.key] = l.Message
}
func (b *buffer) checkLimit() {
	if len(b.queue) > b.limit {
		b.flush()
	}
}

func (b *buffer) flush() {
	for k, o := range b.queue {
		switch o.Message {
		case MessageRemove:
			b.c.next.Remove(k)
		case MessageSet:
			b.c.next.Set(k, o.Value)
		}
	}
	b.reset()
}

// Flusher
func (c *LayerCache) flusher() {
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

			// Next, recursively
			if old.next != nil {
				old.next.Sync()
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
func (c *LayerCache) Get(key interface{}) (value interface{}, err error) {
	// Try to get backend cache
	value, err = c.Storage.Get(key)
	if err != nil {
		if c.next == nil {
			return nil, errors.New("value not found")
		}
		// Recursively get value from list.
		value, err = c.next.Get(key)
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
func (c *LayerCache) Set(key interface{}, value interface{}) (err error) {
	err = c.Storage.Add(key, value)
	if err != nil {
		return err
	}
	if c.next == nil {
		// This backend cache is final destination
		return nil
	}
	// Queue to flush
	c.log <- log{key, &Message{Value: value, Message: MessageSet}}
	return nil
}

// Sync current buffered value
func (c *LayerCache) Sync() error {
	c.sync <- true
	<-c.synced
	return nil
}

// Remove recursively remove next layer's value
func (c *LayerCache) Remove(key interface{}) (err error) {
	err = c.Storage.Remove(key)
	if err != nil {
		return err
	}
	if c.next == nil {
		// This is bottom layer
		return nil
	}
	// Queue to flush
	c.log <- log{key, &Message{Value: nil, Message: MessageRemove}}
	c.Sync() // Remove must be synced
	return nil
}

// SetNext set next layer
func (c *LayerCache) setNext(next Layer) error {
	c.next = next
	return nil
}
