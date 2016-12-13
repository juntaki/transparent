// Package test is simple wrapper of map[interface{}]interface{} for mock
package test

import (
	"sync"
	"time"

	"github.com/juntaki/transparent"
)

type storage struct {
	lock sync.RWMutex
	list map[interface{}]interface{}
	wait time.Duration
}

// NewStorage returns Storage
func NewStorage(wait time.Duration) transparent.Storage {
	return &storage{
		list: make(map[interface{}]interface{}, 0),
		wait: wait,
	}
}

// Get returns value from map
func (d *storage) Get(k interface{}) (interface{}, error) {
	time.Sleep(d.wait * time.Millisecond)
	d.lock.RLock()
	defer d.lock.RUnlock()
	value, ok := d.list[k]
	if !ok {
		return nil, &transparent.KeyNotFoundError{Key: k}
	}
	return value, nil
}

// Add insert value to map
func (d *storage) Add(k interface{}, v interface{}) error {
	time.Sleep(d.wait * time.Millisecond)

	d.lock.Lock()
	defer d.lock.Unlock()
	d.list[k] = v
	return nil
}

// Remove deletes key from map
func (d *storage) Remove(k interface{}) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	delete(d.list, k)
	return nil
}
