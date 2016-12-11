package transparent

import "errors"

// customStorage define customizable storage
type customStorage struct {
	getFunc    func(k interface{}) (interface{}, error)
	addFunc    func(k interface{}, v interface{}) error
	removeFunc func(k interface{}) error
}

// NewCustomStorage returns customStorage
func NewCustomStorage(
	getFunc func(k interface{}) (interface{}, error),
	addFunc func(k interface{}, v interface{}) error,
	removeFunc func(k interface{}) error,
) (Storage, error) {
	if getFunc == nil || addFunc == nil || removeFunc == nil {
		return nil, errors.New("function must be filled")
	}
	return &customStorage{
		getFunc:    getFunc,
		addFunc:    addFunc,
		removeFunc: removeFunc,
	}, nil
}

// Get is customizable get function
func (c *customStorage) Get(k interface{}) (interface{}, error) {
	return c.getFunc(k)
}

// Add is customizable add function
func (c *customStorage) Add(k interface{}, v interface{}) error {
	return c.addFunc(k, v)
}

// Remove is customizable remove function
func (c *customStorage) Remove(k interface{}) error {
	return c.removeFunc(k)
}
