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
	c.SetWriteBack(100, "test")
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
	tiered.SetWriteThrough(100, "test")

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

func BenchmarkTransparentCacheSetWriteBack(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := rand.Intn(5)
		c.SetWriteBack(r, "benchmarking")
	}
}

func BenchmarkTransparentCacheSetWriteThrough(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := rand.Intn(5)
		c.SetWriteThrough(r, "benchmarking")
	}
}

// Tiered
func BenchmarkTieredTransparentCacheGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := rand.Intn(5)
		tiered.Get(r)
	}
}

func BenchmarkTieredTransparentCacheSetWriteBack(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := rand.Intn(5)
		tiered.SetWriteBack(r, "benchmarking")
	}
}

func BenchmarkTieredTransparentCacheSetWriteThrough(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := rand.Intn(5)
		tiered.SetWriteThrough(r, "benchmarking")
	}
}
