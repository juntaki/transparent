package transparent

// Source provides operation of TransparentSource
type Source struct {
	Storage Storage
	upper   Layer
}

// NewSource returns Source
func NewSource() *Source {
	return &Source{}
}

func (s *Source) SetStorage(storage Storage) {
	s.Storage = storage
}

// Set set new value to storage.
func (s *Source) Set(key interface{}, value interface{}) {
	if s.upper != nil {
		s.Skim(key)
	}
	s.Storage.Add(key, value)
}

// Get value from storage
func (s *Source) Get(key interface{}) (value interface{}) {
	value, _ = s.Storage.Get(key)
	return
}

// Remove value
func (s *Source) Remove(key interface{}) {
	s.Storage.Remove(key)
}

// Skim remove upper layer's old value
func (s *Source) Skim(key interface{}) {
	s.Storage.Remove(key)
	if s.upper == nil {
		// This is top layer
		return
	}
	s.upper.Skim(key)
}

// Sync do nothing
func (s *Source) Sync() {
	return
}

func (s *Source) setUpper(upper Layer) {
	s.upper = upper
}

func (s *Source) setLower(lower Layer) {
	panic("don't set lower layer")
}
