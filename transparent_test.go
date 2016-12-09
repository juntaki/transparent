package transparent

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
)

func basicLayerFunc(t *testing.T, l Layer) {
	err := l.Set("test", []byte("value"))
	if err != nil {
		t.Error(err)
	}

	value, err := l.Get("test")
	if err != nil || string(value.([]byte)) != "value" {
		t.Error(err)
		t.Error(value)
	}

	err = l.Remove("test")
	if err != nil {
		t.Error(err)
	}

	value, err = l.Get("test")
	if err == nil {
		t.Error(err)
		t.Error(value)
	}

	err = l.Sync()
	if err != nil {
		t.Error(err)
	}
}

var dummySrc *Source
var dummyCache *Cache

func TestMain(m *testing.M) {
	MyInit()
	retCode := m.Run()
	MyTeardown()
	os.Exit(retCode)
}

func MyInit() {
	var err error
	rand.Seed(time.Now().UnixNano())
	dummySrc, err = NewDummySource(5)
	if err != nil {
		fmt.Println(err)
	}
	dummyCache, err = NewLRUCache(10, 100)
	if err != nil {
		fmt.Println(err)
	}

	Stack(dummyCache, dummySrc)
}

func MyTeardown() {
	DeleteCache(dummyCache)
}

func TestStopFlusher(t *testing.T) {
	src, err := NewDummySource(5)
	if err != nil {
		t.Error(err)
	}
	cache, err := NewLRUCache(10, 100)
	if err != nil {
		t.Error(err)
	}
	Stack(cache, src)
	for i := 0; i < 100; i++ {
		cache.Set(i, i)
	}
	DeleteCache(cache)

	value, err := cache.Get(99)
	if err != nil {
		t.Error(err)
	}
	if value != 99 {
		t.Error(value)
	}
}

func TestSync(t *testing.T) {
	for i := 0; i < 349; i++ {
		dummyCache.Set(i, "writeback")
	}
	dummyCache.Sync()
	value1, err := dummySrc.Get(300)
	if err != nil {
		t.Error(err)
	}
	value2, err := dummyCache.Get(300)
	if err != nil {
		t.Error(err)
	}
	if value1 != value2 {
		t.Error(value1, value2)
	}
}

func BenchmarkSrcGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := rand.Intn(5)
		dummySrc.Get(r)
	}
}

func BenchmarkSrcSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := rand.Intn(5)
		dummySrc.Set(r, "benchmarking")
	}
}

// Tiered
func BenchmarkCacheGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := rand.Intn(5)
		dummyCache.Get(r)
	}
}

func BenchmarkCacheSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := rand.Intn(5)
		dummyCache.Set(r, "benchmarking")
	}
	dummyCache.Sync()
}

func BenchmarkCacheSetAndSync(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := rand.Intn(5)
		dummyCache.Set(r, "benchmarking")
		dummyCache.Sync()
	}
}
