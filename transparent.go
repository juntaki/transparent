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

// Layer is stackable function
type Layer interface {
	Set(key interface{}, value interface{}) error
	Get(key interface{}) (value interface{}, err error)
	Remove(key interface{}) error
	Sync() error
	setUpper(Layer)
	setLower(Layer)
}

// Stack stacks layers
func Stack(upper Layer, lower Layer) {
	upper.setLower(lower)
	lower.setUpper(upper)
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
