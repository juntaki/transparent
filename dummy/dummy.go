package dummy

import (
	"sync"
	"time"

	"github.com/juntaki/transparent"
)

// dummyStorage is simple wrapper of map[interface{}]interface{} for mock
type storage struct {
	lock sync.RWMutex
	list map[interface{}]interface{}
	wait time.Duration
}

// NewDummyStorage returns dummy storage
func NewStorage(wait time.Duration) (transparent.Storage, error) {
	return &storage{
		list: make(map[interface{}]interface{}, 0),
		wait: wait,
	}, nil
}

// Get returns value from map
func (d *storage) Get(k interface{}) (interface{}, error) {
	time.Sleep(d.wait * time.Millisecond)
	d.lock.RLock()
	defer d.lock.RUnlock()
	value, ok := d.list[k]
	if !ok {
		return nil, &transparent.StorageKeyNotFoundError{Key: k}
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
