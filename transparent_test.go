package transparent

import (
	"math/rand"
	"os"
	"testing"
	"time"
)

var dummySrc *Source
var dummyCache *Cache

func TestMain(m *testing.M) {
	MyInit()
	retCode := m.Run()
	MyTeardown()
	os.Exit(retCode)
}

func MyInit() {
	rand.Seed(time.Now().UnixNano())
	dummySrc = NewDummySource(5)
	dummyCache = NewLRUCache(100, 10)
	Stack(dummyCache, dummySrc)
}

func MyTeardown() {
	DeleteCache(dummyCache)
}

func TestStopFlusher(t *testing.T) {
	src := NewDummySource(5)
	cache := NewLRUCache(100, 10)
	Stack(cache, src)
	for i := 0; i < 100; i++ {
		cache.Set(i, i)
	}
	DeleteCache(cache)

	value := cache.Get(99)
	if value != 99 {
		t.Error(value)
	}
}

// Simple Set and Get
func TestSrc(t *testing.T) {
	dummySrc.Set(100, "test")
	value := dummySrc.Get(100)
	if value != "test" {
		t.Error(value)
	}
	dummySrc.Remove(100)
}

// Tiered, Set and Get
func TestCache(t *testing.T) {
	dummySrc.Set(100, "test")
	value := dummyCache.Get(100)
	if value != "test" {
		t.Error(value)
	}
	dummyCache.Set(100, "test")
	dummyCache.Sync()

	value = dummyCache.Get(100)
	if value != "test" {
		t.Error(value)
	}
	dummyCache.Remove(100)
}

func TestSync(t *testing.T) {
	for i := 0; i < 349; i++ {
		dummyCache.Set(i, "writeback")
	}
	dummyCache.Sync()
	value1 := dummySrc.Get(300)
	value2 := dummyCache.Get(300)
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
