package test

import (
	"testing"

	"github.com/juntaki/transparent"
)

func TestDummy(t *testing.T) {
	ds := NewStorage(0)
	BasicStorageFunc(t, ds)

	l := NewSource(1)
	s := transparent.NewStack()
	s.Stack(l)
	BasicStackFunc(t, s)
}
