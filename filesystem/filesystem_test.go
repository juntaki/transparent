package filesystem

import (
	"testing"

	"github.com/juntaki/transparent"
	test "github.com/juntaki/transparent/test"
)

func TestFilesystemSimpleStorage(t *testing.T) {
	fss := NewSimpleStorage("/tmp")
	test.BasicStorageFunc(t, fss)
	test.SimpleStorageFunc(t, fss)
}

func TestFilesystemStorage(t *testing.T) {
	fs := NewStorage("/tmp")
	test.BasicStorageFunc(t, fs)
}

func TestFilesystemSource(t *testing.T) {
	l := NewSource("/tmp")
	stack := transparent.NewStack()
	stack.Stack(l)
	test.BasicStackFunc(t, stack)
}

func TestFilesystemCache(t *testing.T) {
	c := NewCache(10, "/tmp")
	test.BasicCacheFunc(t, c)
}
