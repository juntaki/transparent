package test

import (
	"time"

	"github.com/juntaki/transparent"
)

// NewSource returns Source
func NewSource(wait time.Duration) *transparent.LayerSource {
	test := NewStorage(wait)
	layer, _ := transparent.NewLayerSource(test)
	return layer
}
