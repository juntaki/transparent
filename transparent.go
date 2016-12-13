// Package transparent is a library that provides transparent caching operations
// for key-value stores. As shown in the figure below, it is possible to use relatively
// fast cache like LRU and slow and reliable storage like S3 via TransparentCache.
// Transparent Cache is tearable. In addition to caching, it is also possible to
// transparently use a layer of synchronization between distributed systems.
// See subpackage for example.

package transparent

// Stack is stacked layer
type Stack struct {
	Layer
	all []Layer
}

// NewStack returns Stack
func NewStack() *Stack {
	return &Stack{
		all: []Layer{},
	}
}

// Stack add the layer to Stack
func (s *Stack) Stack(l Layer) error {
	if s.Layer != nil {
		err := l.setLower(s.Layer)
		if err != nil {
			return err
		}
		err = s.Layer.setUpper(l)
		if err != nil {
			return err
		}
	}
	s.Layer = l
	s.all = append(s.all, l)
	return nil
}

// Start initialize all stacked layers
func (s *Stack) Start() error {
	for _, l := range s.all {
		err := l.start()
		if err != nil {
			return err
		}
	}
	return nil
}

// Stop clean up all stacked layers
func (s *Stack) Stop() error {
	for _, l := range s.all {
		err := l.stop()
		if err != nil {
			return err
		}
	}
	return nil
}

// Layer is stackable function
type Layer interface {
	Set(key interface{}, value interface{}) error
	Get(key interface{}) (value interface{}, err error)
	Remove(key interface{}) error
	Sync() error
	setUpper(Layer) error
	setLower(Layer) error
	start() error
	stop() error
}

// message passing between layer or its internals
type message int

const (
	messageSet message = iota
	messageRemove
	messageSync
)

type operation struct {
	Value   interface{}
	Message message
	UUID    string
}

// KeyNotFoundError means specified key is not found in the layer
type KeyNotFoundError struct {
	Key interface{}
}

func (e *KeyNotFoundError) Error() string { return "requested key is not found" }
