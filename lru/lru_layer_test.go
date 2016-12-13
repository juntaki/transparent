package lru

import (
	"testing"

	test "github.com/juntaki/transparent/test"
)

func TestLRUCache(t *testing.T) {
	c, err := NewCache(10, 100)
	if err != nil {
		t.Error(err)
	}
	test.BasicCacheFunc(t, c)
}
