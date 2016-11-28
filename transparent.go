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
//   `-[Backend cache]                     `-[Source cache]
//      `-[LRU]                               `-[S3]
package transparent

// Layer is stackable function
type Layer interface {
	Set(key interface{}, value interface{})
	Get(key interface{}) (value interface{})
	Remove(key interface{})
	Skim(key interface{})
	Sync()
}

type stackable interface {
	setUpper(Layer)
	getLayer() Layer
	setLower(Layer)
}

// stacker implements stackable interface
type stacker struct {
	upper Layer
	this  Layer
	lower Layer
}

func (s *stacker) setUpper(l Layer) {
	s.upper = l
}
func (s *stacker) getLayer() Layer {
	return s.this
}
func (s *stacker) setLower(l Layer) {
	s.lower = l
}

func (s *stacker) Stack(upper stackable) {
	s.setUpper(upper.getLayer())
	upper.setLower(s.getLayer())
}
