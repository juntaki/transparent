package custom

import (
	"errors"

	"github.com/juntaki/transparent"
)

type storage struct {
	getFunc    func(k interface{}) (interface{}, error)
	addFunc    func(k interface{}, v interface{}) error
	removeFunc func(k interface{}) error
}

// NewStorage returns Storage
func NewStorage(
	getFunc func(k interface{}) (interface{}, error),
	addFunc func(k interface{}, v interface{}) error,
	removeFunc func(k interface{}) error,
) (transparent.Storage, error) {
	if getFunc == nil || addFunc == nil || removeFunc == nil {
		return nil, errors.New("function must be filled")
	}
	return &storage{
		getFunc:    getFunc,
		addFunc:    addFunc,
		removeFunc: removeFunc,
	}, nil
}

// Get is customizable get function
func (c *storage) Get(k interface{}) (interface{}, error) {
	return c.getFunc(k)
}

// Add is customizable add function
func (c *storage) Add(k interface{}, v interface{}) error {
	return c.addFunc(k, v)
}

// Remove is customizable remove function
func (c *storage) Remove(k interface{}) error {
	return c.removeFunc(k)
}
