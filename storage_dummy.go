package transparent

import (
	"sync"
	"time"
)

// DummyStorage is simple wrapper of map[interface{}]interface{} for mock
type DummyStorage struct {
	lock sync.RWMutex
	list map[interface{}]interface{}
	wait time.Duration
}

// NewDummyStorage returns dummy storage
func NewDummyStorage(wait time.Duration) (*DummyStorage, error) {
	return &DummyStorage{
		list: make(map[interface{}]interface{}, 0),
		wait: wait,
	}, nil
}

// Get returns value from map
func (d *DummyStorage) Get(k interface{}) (interface{}, error) {
	time.Sleep(d.wait * time.Millisecond)
	d.lock.RLock()
	defer d.lock.RUnlock()
	value, ok := d.list[k]
	if !ok {
		return nil, &StorageKeyNotFoundError{Key: k}
	}
	return value, nil
}

// Add insert value to map
func (d *DummyStorage) Add(k interface{}, v interface{}) error {
	time.Sleep(d.wait * time.Millisecond)

	d.lock.Lock()
	defer d.lock.Unlock()
	d.list[k] = v
	return nil
}

// Remove deletes key from map
func (d *DummyStorage) Remove(k interface{}) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	delete(d.list, k)
	return nil
}
