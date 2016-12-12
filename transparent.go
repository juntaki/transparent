// Package transparent is a library that provides transparent caching operations
// for key-value stores. As shown in the figure below, it is possible to use relatively
// fast cache like LRU and slow and reliable storage like S3 via TransparentCache.
// Transparent Cache is tearable. In addition to caching, it is also possible to
// transparently use a layer of synchronization between distributed systems.
// See subpackage for example.
//
//  [Application]
//    |
//    v Get/Set
//  [Transparent cache] -[Flush buffer]-> [Next cache]
//   `-[backend cache]                     `-[Source cache]
//      `-[LRU]                               `-[S3]
package transparent

type Stack struct {
	Layer
	all []Layer
}

func NewStack() *Stack {
	return &Stack{
		all: []Layer{},
	}
}

func (s *Stack) Stack(l Layer) error {
	if s.Layer != nil {
		l.setLower(s.Layer)
		s.Layer.setUpper(l)
	}
	s.Layer = l
	s.all = append(s.all, l)
	return nil
}

func (s *Stack) Start() error {
	for _, l := range s.all {
		err := l.start()
		if err != nil {
			return err
		}
	}
	return nil
}

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
	setUpper(Layer)
	setLower(Layer)
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
