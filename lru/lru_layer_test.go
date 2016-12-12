package lru

import (
	"testing"

	"github.com/juntaki/transparent"
	test "github.com/juntaki/transparent/test"
)

func TestCache(t *testing.T) {
	var c *transparent.Cache
	var err error
	c, err = NewCache(10, 100)
	if err != nil {
		t.Error(err)
	}
	test.BasicCacheFunc(t, c)
}
