package lru

import (
	"github.com/juntaki/transparent"
)

// NewCache returns LRUCache
func NewCache(bufferSize, cacheSize int) (*transparent.Cache, error) {
	lru := New(cacheSize)
	layer, err := transparent.NewCacheLayer(bufferSize, lru)
	if err != nil {
		return nil, err
	}
	return layer, nil
}
