package transparent

import (
	"fmt"
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
	var err error
	rand.Seed(time.Now().UnixNano())
	dummySrc = NewDummySource(5)
	dummyCache, err = NewLRUCache(10, 100)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(dummyCache.Storage)
	Stack(dummyCache, dummySrc)
}

func MyTeardown() {
	DeleteCache(dummyCache)
}

func TestStopFlusher(t *testing.T) {
	src := NewDummySource(5)
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

// Simple Set and Get
func TestSrc(t *testing.T) {
	dummySrc.Set(100, "test")
	value, err := dummySrc.Get(100)
	if err != nil {
		t.Error(err)
	}
	if value != "test" {
		t.Error(value)
	}
	dummySrc.Remove(100)
}

// Tiered, Set and Get
func TestCache(t *testing.T) {
	dummySrc.Set(100, "test")
	value, err := dummyCache.Get(100)
	if err != nil {
		t.Error(err)
	}
	if value != "test" {
		t.Error(value)
	}
	dummyCache.Set(100, "test")
	dummyCache.Sync()

	value, err = dummyCache.Get(100)
	if err != nil {
		t.Error(err)
	}
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
