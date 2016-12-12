package dummy

import (
	"time"

	"github.com/juntaki/transparent"
)

// NewDummySource returns dummyStorage layer
func NewSource(wait time.Duration) (*transparent.Source, error) {
	dummy, err := NewStorage(wait)
	if err != nil {
		return nil, err
	}
	layer, _ := transparent.NewSource(dummy)
	return layer, nil
}
