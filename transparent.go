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

// Layer
type Layer interface {
	Set(key interface{}, value interface{})
	Get(key interface{}) (value interface{})
	Remove(key interface{})
	Sync()
	Stack(Layer)
	getListHead() *listHead
}

type listHead struct {
	upper Layer
	lower Layer
}
