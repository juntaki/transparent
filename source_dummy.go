package transparent

import (
	"sync"
	"time"
)

// Define dummy source
type dummySource struct {
	lock sync.RWMutex
	list map[interface{}]interface{}
	wait time.Duration
}

func (d *dummySource) Get(k interface{}) (interface{}, error) {
	time.Sleep(d.wait * time.Millisecond)
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.list[k], nil
}
func (d *dummySource) Add(k interface{}, v interface{}) error {
	time.Sleep(d.wait * time.Millisecond)

	d.lock.Lock()
	defer d.lock.Unlock()
	d.list[k] = v
	return nil
}
func (d *dummySource) Remove(k interface{}) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	delete(d.list, k)
	return nil
}

// NewDummySource returns dummySource layer
func NewDummySource(wait time.Duration) *Source {
	layer, _ := NewSource(&dummySource{
		list: make(map[interface{}]interface{}, 0),
		wait: wait,
	})
	return layer
}
