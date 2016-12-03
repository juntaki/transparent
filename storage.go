package transparent

import "errors"

// Storage defines the interface that backend data storage destination should have.
// Add should not be failed.
type Storage interface {
	Get(key interface{}) (value interface{}, err error)
	Add(key interface{}, value interface{}) error
	Remove(key interface{}) error
}

// CustomStorage define customizable storage
type CustomStorage struct {
	getFunc    func(k interface{}) (interface{}, error)
	addFunc    func(k interface{}, v interface{}) error
	removeFunc func(k interface{}) error
}

// Get is customizable get function
func (c *CustomStorage) Get(k interface{}) (interface{}, error) {
	return c.getFunc(k)
}

// Add is customizable add function
func (c *CustomStorage) Add(k interface{}, v interface{}) error {
	return c.addFunc(k, v)
}

// Remove is customizable remove function
func (c *CustomStorage) Remove(k interface{}) error {
	return c.removeFunc(k)
}

// NewCustomStorage returns CustomStorage
func NewCustomStorage(
	getFunc func(k interface{}) (interface{}, error),
	addFunc func(k interface{}, v interface{}) error,
	removeFunc func(k interface{}) error,
) (*CustomStorage, error) {
	if getFunc == nil || addFunc == nil || removeFunc == nil {
		return nil, errors.New("function must be filled")
	}
	return &CustomStorage{
		getFunc:    getFunc,
		addFunc:    addFunc,
		removeFunc: removeFunc,
	}, nil
}
