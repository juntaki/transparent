package test

import (
	"time"

	"github.com/juntaki/transparent"
)

// NewSource returns Source
func NewSource(wait time.Duration) *transparent.Source {
	test := NewStorage(wait)
	layer, _ := transparent.NewSource(test)
	return layer
}
