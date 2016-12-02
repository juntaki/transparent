package transparent

import lru "github.com/juntaki/transparent/lru"

// NewLRUCache returns transparent cache
func NewLRUCache(cacheSize, bufferSize int) *Cache {
	layer := New(bufferSize)
	layer.backendCache = lru.New(cacheSize)
	layer.startFlusher()
	return layer
}
