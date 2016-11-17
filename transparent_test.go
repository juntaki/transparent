package transparent

import (
	"math/rand"
	"testing"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

type dummySource struct {
	list map[int]string
}

func (d dummySource) Get(k interface{}) (interface{}, bool) {
	time.Sleep(100 * time.Millisecond)

	return d.list[k.(int)], true
}
func (d dummySource) Add(k, v interface{}) bool {
	time.Sleep(100 * time.Millisecond)
	d.list[k.(int)] = v.(string)
	return true
}

var d dummySource
var c Cache
var tiered Cache

func init() {
	rand.Seed(time.Now().UnixNano())
	d = dummySource{}
	d.list = make(map[int]string, 0)
	c = Cache{
		cache: d,
		next:  nil,
	}

	lru, err := lru.New(10)
	if err != nil {
		panic("LRU error")
	}
	tiered = Cache{
		cache: lru,
		next:  &c,
	}
}

func TestTransparentCache(t *testing.T) {
	c.Set(100, "test")
	value := c.Get(100)
	if value != "test" {
		t.Error(value)
	}
}

func TestTieredTransparentCache(t *testing.T) {
	value := tiered.Get(100)
	if value != "test" {
		t.Error(value)
	}
	tiered.SetSync(100, "test")

	value = tiered.Get(100)
	if value != "test" {
		t.Error(value)
	}
}

func BenchmarkTransparentCacheGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := rand.Intn(5)
		c.Get(r)
	}
}

func BenchmarkTransparentCacheSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := rand.Intn(5)
		c.Set(r, "benchmarking")
	}
}

func BenchmarkTransparentCacheSetSync(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := rand.Intn(5)
		c.SetSync(r, "benchmarking")
	}
}

// Tiered
func BenchmarkTieredTransparentCacheGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := rand.Intn(5)
		tiered.Get(r)
	}
}

func BenchmarkTieredTransparentCacheSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := rand.Intn(5)
		tiered.Set(r, "benchmarking")
	}
}

func BenchmarkTieredTransparentCacheSetSync(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := rand.Intn(5)
		tiered.SetSync(r, "benchmarking")
	}
}
