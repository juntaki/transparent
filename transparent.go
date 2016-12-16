// Package transparent is a library that provides transparent operations for key-value stores.
// Transparent Layer is tearable on Stack. In addition to caching, it is also possible to
// transparently use a layer of synchronization between distributed systems.
// See subpackage for implementation.
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
		err := l.setNext(s.Layer)
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
	setNext(Layer) error
	start() error
	stop() error
}

// MessageType of operation
type MessageType int

// MessageType of operation
const (
	MessageSet MessageType = iota
	MessageGet
	MessageRemove
	MessageSync
)

// Message is layer operation
type Message struct {
	Key     interface{}
	Value   interface{}
	Message MessageType
	UUID    string
}

// KeyNotFoundError means specified key is not found in the layer
type KeyNotFoundError struct {
	Key interface{}
}

func (e *KeyNotFoundError) Error() string { return "requested key is not found" }
