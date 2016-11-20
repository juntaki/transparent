package transparent

import (
	"math/rand"
	"os"
	"testing"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

func TestMain(m *testing.M) {
	MyInit()
	retCode := m.Run()
	MyTeardown()
	os.Exit(retCode)
}

// Define dummy source
type dummySource struct {
	list map[int]string
}

func (d dummySource) Get(k interface{}) (interface{}, bool) {
	time.Sleep(5 * time.Millisecond)

	return d.list[k.(int)], true
}
func (d dummySource) Add(k, v interface{}) bool {
	time.Sleep(5 * time.Millisecond)
	d.list[k.(int)] = v.(string)
	return true
}

var d dummySource
var c Cache
var tiered Cache

func MyInit() {
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
	c.Initialize(300)
	tiered.Initialize(300)
}

func MyTeardown() {
	c.Finalize()
	tiered.Finalize()
}

func TestFinalize(t *testing.T) {
	lru, err := lru.New(10)
	if err != nil {
		panic("LRU error")
	}
	cache := Cache{
		cache: lru,
		next:  &c,
	}
	cache.Initialize(100)
	cache.SetWriteThrough(100, "Test")
	cache.Finalize()
}

// Simple Set and Get
func TestCache(t *testing.T) {
	c.SetWriteBack(100, "test")
	value := c.Get(100)
	if value != "test" {
		t.Error(value)
	}
}

// Tiered, Set and Get
func TestTieredCache(t *testing.T) {
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

func BenchmarkCacheGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := rand.Intn(5)
		c.Get(r)
	}
}

func BenchmarkCacheSetWriteBack(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := rand.Intn(5)
		c.SetWriteBack(r, "benchmarking")
	}
}

func BenchmarkCacheSetWriteThrough(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := rand.Intn(5)
		c.SetWriteThrough(r, "benchmarking")
	}
}

// Tiered
func BenchmarkTieredCacheGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := rand.Intn(5)
		tiered.Get(r)
	}
}

func BenchmarkTieredCacheSetWriteBack(b *testing.B) {
	for i := 0; i < 100; i++ {
		r := rand.Intn(5)
		tiered.SetWriteBack(r, "benchmarking")
	}
}

func BenchmarkTieredCacheSetWriteThrough(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := rand.Intn(5)
		tiered.SetWriteThrough(r, "benchmarking")
	}
}
