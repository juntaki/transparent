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
	if s.upper != nil {
		err = s.Skim(key)
		if err != nil {
			return err
		}
	}
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

// Skim remove upper layer's old value
func (s *Source) Skim(key interface{}) (err error) {
	err = s.Storage.Remove(key)
	if err != nil {
		return err
	}
	if s.upper == nil {
		// This is top layer
		return
	}
	err = s.upper.Skim(key)
	if err != nil {
		return err
	}
	return nil
}

// Sync do nothing
func (s *Source) Sync() error {
	return nil
}

func (s *Source) setUpper(upper Layer) {
	s.upper = upper
}

func (s *Source) setLower(lower Layer) {
	panic("don't set lower layer")
}
