package transparent

import "errors"

type layerSource struct {
	Storage BackendStorage
}

// NewLayerSource returns LayerSource.
// LayerSource wraps BackendStorage.
// It Get/Set key-value to BackendStorage.
// This layer must be the bottom of Stack.
func NewLayerSource(storage BackendStorage) (Layer, error) {
	if storage == nil {
		return nil, errors.New("empty storage")
	}
	return &layerSource{Storage: storage}, nil
}

// Set set new value to storage.
func (s *layerSource) Set(key interface{}, value interface{}) (err error) {
	err = s.Storage.Add(key, value)
	if err != nil {
		return err
	}
	return nil
}

// Get value from storage
func (s *layerSource) Get(key interface{}) (value interface{}, err error) {
	return s.Storage.Get(key)
}

// Remove value
func (s *layerSource) Remove(key interface{}) (err error) {
	return s.Storage.Remove(key)
}

// Sync do nothing
func (s *layerSource) Sync() error {
	return nil
}

func (s *layerSource) setNext(next Layer) error {
	return errors.New("don't set next layer")
}

func (s *layerSource) start() error {
	return nil
}

func (s *layerSource) stop() error {
	return nil
}
