package transparent

// Storage defines the interface that backend data storage destination should have.
// Add should not be failed.
type Storage interface {
	Get(key interface{}) (value interface{}, found bool)
	Add(key interface{}, value interface{})
	Remove(key interface{})
}

// CustomStorage define customizable storage
type CustomStorage struct {
	get    func(k interface{}) (interface{}, bool)
	add    func(k interface{}, v interface{})
	remove func(k interface{})
}

func (c *CustomStorage) Get(k interface{}) (interface{}, bool) {
	return c.get(k)
}
func (c *CustomStorage) Add(k interface{}, v interface{}) {
	c.add(k, v)
}
func (c *CustomStorage) Remove(k interface{}) {
	c.remove(k)
}
