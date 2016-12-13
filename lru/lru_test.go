package lru

import (
	"testing"

	test "github.com/juntaki/transparent/test"
)

func TestLRUStorage(t *testing.T) {
	c := NewStorage(10)
	test.BasicStorageFunc(t, c)
}
