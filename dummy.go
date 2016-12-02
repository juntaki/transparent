package transparent

import "time"

// Define dummy source
type dummySource struct {
	list map[interface{}]interface{}
	wait time.Duration
}

func (d *dummySource) Get(k interface{}) (interface{}, bool) {
	time.Sleep(d.wait * time.Millisecond)
	return d.list[k], true
}
func (d *dummySource) Add(k interface{}, v interface{}) {
	time.Sleep(d.wait * time.Millisecond)
	d.list[k] = v
}
func (d *dummySource) Remove(k interface{}) {
	delete(d.list, k)
}

// NewDummySource returns dummySource layer
func NewDummySource(wait time.Duration) *Source {
	layer := NewSource()
	layer.storage = &dummySource{
		list: make(map[interface{}]interface{}, 0),
		wait: wait,
	}
	return layer
}
