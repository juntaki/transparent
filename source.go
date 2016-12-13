package transparent

import "errors"

// Source provides operation of TransparentSource
type Source struct {
	Storage Storage
	upper   Layer
}

// NewSource returns Source
func NewSource(storage Storage) (*Source, error) {
	if storage == nil {
		return nil, errors.New("empty storage")
	}
	return &Source{Storage: storage}, nil
}

// Set set new value to storage.
func (s *Source) Set(key interface{}, value interface{}) (err error) {
	err = s.Storage.Add(key, value)
	if err != nil {
		return err
	}
	return nil
}

// Get value from storage
func (s *Source) Get(key interface{}) (value interface{}, err error) {
	return s.Storage.Get(key)
}

// Remove value
func (s *Source) Remove(key interface{}) (err error) {
	return s.Storage.Remove(key)
}

// Sync do nothing
func (s *Source) Sync() error {
	return nil
}

func (s *Source) setUpper(upper Layer) error {
	s.upper = upper
	return nil
}

func (s *Source) setLower(lower Layer) error {
	return errors.New("don't set lower layer")
}

func (s *Source) start() error {
	return nil
}

func (s *Source) stop() error {
	return nil
}
