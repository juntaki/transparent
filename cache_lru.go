package transparent

import lru "github.com/juntaki/transparent/lru"

// NewLRUCache returns transparent cache
func NewLRUCache(cacheSize, bufferSize int) *Cache {
	layer := NewCache(bufferSize)
	layer.storage = lru.New(cacheSize)
	layer.startFlusher()
	return layer
}
