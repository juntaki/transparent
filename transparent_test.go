package transparent

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

type dummySource struct{}

func (d dummySource) Get(k interface{}) (interface{}, bool) {
	time.Sleep(100 * time.Millisecond)
	return "test", true
}
func (d dummySource) Add(k, v interface{}) bool {
	time.Sleep(100 * time.Millisecond)
	return true
}

func TestTransparentCache(t *testing.T) {
	var d dummySource
	c := Cache{
		cache: d,
		next:  nil,
	}

	value, _ := c.Get(100)
	fmt.Println(value)
	c.Set(100, 100)
	value, _ = c.Get(100)
	fmt.Println(value)
}

func TestTieredTransparentCache(t *testing.T) {
	var d dummySource
	source := Cache{
		cache: d,
		next:  nil,
	}

	lru, err := lru.New(10)
	if err != nil {
		panic("LRU error")
	}
	c := Cache{
		cache: lru,
		next:  source,
	}

	value, _ := c.Get(100)
	fmt.Println(value)
	c.Set(100, 100)
	value, _ = c.Get(100)
	fmt.Println(value)
}

func BenchmarkTransparentCacheGet(b *testing.B) {
	var d dummySource
	c := Cache{
		cache: d,
		next:  nil,
	}

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < b.N; i++ {
		r := rand.Intn(5)
		c.Get(r)
		c.Set(r, "benchmarking")
	}
}

func BenchmarkTieredTransparentCacheGet(b *testing.B) {
	var d dummySource
	source := Cache{
		cache: d,
		next:  nil,
	}

	lru, err := lru.New(10)
	if err != nil {
		panic("LRU error")
	}
	c := Cache{
		cache: lru,
		next:  source,
	}

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < b.N; i++ {
		r := rand.Intn(5)
		c.Get(r)
		c.Set(r, "benchmarking")
	}
}
