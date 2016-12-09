package transparent

import (
	"errors"
	"sync"
	"time"
)

// Define dummy storage
type DummyStorage struct {
	lock sync.RWMutex
	list map[interface{}]interface{}
	wait time.Duration
}

func NewDummyStorage(wait time.Duration) (*DummyStorage, error) {
	return &DummyStorage{
		list: make(map[interface{}]interface{}, 0),
		wait: wait,
	}, nil
}

func (d *DummyStorage) Get(k interface{}) (interface{}, error) {
	time.Sleep(d.wait * time.Millisecond)
	d.lock.RLock()
	defer d.lock.RUnlock()
	value, ok := d.list[k]
	if !ok {
		return nil, errors.New("key not found")
	}
	return value, nil
}
func (d *DummyStorage) Add(k interface{}, v interface{}) error {
	time.Sleep(d.wait * time.Millisecond)

	d.lock.Lock()
	defer d.lock.Unlock()
	d.list[k] = v
	return nil
}
func (d *DummyStorage) Remove(k interface{}) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	delete(d.list, k)
	return nil
}
