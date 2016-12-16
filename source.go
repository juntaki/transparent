package transparent

import "errors"

// LayerSource wraps BackendStorage.
// It Get/Set key-value to BackendStorage.
// This layer must be the bottom of Stack
type LayerSource struct {
	Storage BackendStorage
}

// NewLayerSource returns Source
func NewLayerSource(storage BackendStorage) (*LayerSource, error) {
	if storage == nil {
		return nil, errors.New("empty storage")
	}
	return &LayerSource{Storage: storage}, nil
}

// Set set new value to storage.
func (s *LayerSource) Set(key interface{}, value interface{}) (err error) {
	err = s.Storage.Add(key, value)
	if err != nil {
		return err
	}
	return nil
}

// Get value from storage
func (s *LayerSource) Get(key interface{}) (value interface{}, err error) {
	return s.Storage.Get(key)
}

// Remove value
func (s *LayerSource) Remove(key interface{}) (err error) {
	return s.Storage.Remove(key)
}

// Sync do nothing
func (s *LayerSource) Sync() error {
	return nil
}

func (s *LayerSource) setNext(next Layer) error {
	return errors.New("don't set next layer")
}

func (s *LayerSource) start() error {
	return nil
}

func (s *LayerSource) stop() error {
	return nil
}
