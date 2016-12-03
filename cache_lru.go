package transparent

import lru "github.com/juntaki/transparent/lru"

// NewLRUCache returns transparent cache
func NewLRUCache(cacheSize, bufferSize int) *Cache {
	layer, _ := NewCache(bufferSize, lru.New(cacheSize))
	layer.startFlusher()
	return layer
}
