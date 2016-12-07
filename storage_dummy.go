package transparent

import (
	"sync"
	"time"
)

// Define dummy storage
type dummyStorage struct {
	lock sync.RWMutex
	list map[interface{}]interface{}
	wait time.Duration
}

func (d *dummyStorage) Get(k interface{}) (interface{}, error) {
	time.Sleep(d.wait * time.Millisecond)
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.list[k], nil
}
func (d *dummyStorage) Add(k interface{}, v interface{}) error {
	time.Sleep(d.wait * time.Millisecond)

	d.lock.Lock()
	defer d.lock.Unlock()
	d.list[k] = v
	return nil
}
func (d *dummyStorage) Remove(k interface{}) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	delete(d.list, k)
	return nil
}
