package lru

import lru "github.com/golang/groupcache/lru"

// GroupcacheLRU is google groupcache imprementation.
type GroupcacheLRU struct {
	cache *lru.Cache
}

// Get wraps its function, it's compatible to transparent.BackendCache
func (d *GroupcacheLRU) Get(k interface{}) (interface{}, bool) {
	return d.cache.Get(k)
}

// Add wraps its function, it's compatible to transparent.BackendCache
func (d *GroupcacheLRU) Add(k interface{}, v interface{}) {
	d.cache.Add(k, v)
}

func GroupcacheLRUNew(maxEntries int) *GroupcacheLRU {
	return &GroupcacheLRU{
		cache: lru.New(maxEntries),
	}
}
