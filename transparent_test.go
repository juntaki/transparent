package transparent

import (
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	gl "github.com/golang/groupcache/lru"
	hl "github.com/hashicorp/golang-lru"
	tl "github.com/juntaki/transparent/lru"
)

// Define tiered cache
type hlTier struct {
	cache *hl.Cache
}

func (d hlTier) Get(k interface{}) (interface{}, bool) {
	return d.cache.Get(k)
}
func (d hlTier) Add(k interface{}, v interface{}) {
	d.cache.Add(k, v)
}

type glTier struct {
	cache *gl.Cache
}

func (d glTier) Get(k interface{}) (interface{}, bool) {
	return d.cache.Get(k)
}
func (d glTier) Add(k interface{}, v interface{}) {
	d.cache.Add(k, v)
}

func TestMain(m *testing.M) {
	MyInit(tl.New(100))
	retCode := m.Run()
	MyTeardown()

	hlru, err := hl.New(100)
	if err != nil {
		panic("LRU error")
	}
	tieredCacheHashicorp := hlTier{cache: hlru}
	MyInit(tieredCacheHashicorp)
	retCode = m.Run()
	MyTeardown()

	glru := gl.New(100)
	tieredCacheGoogle := glTier{cache: glru}
	MyInit(tieredCacheGoogle)
	retCode = m.Run()
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
func (d dummySource) Add(k interface{}, v interface{}) {
	time.Sleep(5 * time.Millisecond)
	d.list[k.(int)] = v.(string)
}

var dummybackend0 dummySource
var dummycache0 *Cache
var tieredbackend1 BackendCache
var tieredcache1 *Cache

func MyInit(backend BackendCache) {
	rand.Seed(time.Now().UnixNano())
	dummybackend0 = dummySource{}
	dummybackend0.list = make(map[int]string, 0)
	dummycache0 = New(
		dummybackend0,
		nil,
		300,
	)
	tieredbackend1 = backend
	tieredcache1 = New(
		backend,
		dummycache0,
		300,
	)
	dummycache0.Initialize()
	tieredcache1.Initialize()
}

func MyTeardown() {
	dummycache0.Finalize()
	tieredcache1.Finalize()
}

func TestFinalize(t *testing.T) {
	cache := New(
		tieredbackend1,
		dummycache0,
		100,
	)
	cache.Initialize()
	for i := 0; i < 100; i++ {
		cache.SetWriteBack(i, strconv.Itoa(i))
	}
	cache.Finalize()

	// not proper use
	value := cache.Get(99)
	if value != "99" {
		t.Error(value)
	}
}

// Simple Set and Get
func TestCache(t *testing.T) {
	dummycache0.SetWriteBack(100, "test")
	value := dummycache0.Get(100)
	if value != "test" {
		t.Error(value)
	}
}

// Tiered, Set and Get
func TestTieredCache(t *testing.T) {
	value := tieredcache1.Get(100)
	if value != "test" {
		t.Error(value)
	}
	tieredcache1.SetWriteThrough(100, "test")

	value = tieredcache1.Get(100)
	if value != "test" {
		t.Error(value)
	}
}

func TestConcurrentUpdate(t *testing.T) {
	for i := 0; i < 350; i++ {
		tieredcache1.SetWriteBack(100, "writeback")
	}
	tieredcache1.SetWriteThrough(100, "writethrough")
	value1 := dummycache0.Get(100)
	value2 := tieredcache1.Get(100)
	if value1 != value2 {
		t.Error(value1, value2)
	}
}

func TestSync(t *testing.T) {
	for i := 0; i < 349; i++ {
		tieredcache1.SetWriteBack(i, "writeback")
	}
	tieredcache1.Sync()
	value1 := dummycache0.Get(300)
	value2 := tieredcache1.Get(300)
	if value1 != value2 {
		t.Error(value1, value2)
	}
}

func BenchmarkCacheGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := rand.Intn(5)
		dummycache0.Get(r)
	}
}

func BenchmarkCacheSetWriteBack(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := rand.Intn(5)
		dummycache0.SetWriteBack(r, "benchmarking")
	}
}

func BenchmarkCacheSetWriteThrough(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := rand.Intn(5)
		dummycache0.SetWriteThrough(r, "benchmarking")
	}
}

// Tiered
func BenchmarkTieredCacheGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := rand.Intn(5)
		tieredcache1.Get(r)
	}
}

func BenchmarkTieredCacheSetWriteBack(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := rand.Intn(5)
		tieredcache1.SetWriteBack(r, "benchmarking")
	}
}

func BenchmarkTieredCacheSetWriteThrough(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := rand.Intn(5)
		tieredcache1.SetWriteThrough(r, "benchmarking")
	}
}
