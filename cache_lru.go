package transparent

import lru "github.com/juntaki/transparent/lru"

// NewLRUCache returns transparent cache
func NewLRUCache(cacheSize, bufferSize int) *Cache {
	layer := New(bufferSize)
	layer.BackendCache = lru.New(cacheSize)
	layer.StartFlusher()
	return layer
}
