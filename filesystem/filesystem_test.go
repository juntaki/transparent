package filesystem

import (
	"testing"

	"github.com/juntaki/transparent"
	test "github.com/juntaki/transparent/test"
)

func TestFilesystemSimpleStorage(t *testing.T) {
	fss, err := NewSimpleStorage("/tmp")
	if err != nil {
		t.Fatal(err)
	}
	test.BasicStorageFunc(t, fss)
	test.SimpleStorageFunc(t, fss)
}

func TestFilesystemStorage(t *testing.T) {
	fs, err := NewStorage("/tmp")
	if err != nil {
		t.Fatal(err)
	}
	test.BasicStorageFunc(t, fs)
}

func TestFilesystemSource(t *testing.T) {
	l, err := NewSource("/tmp")
	if err != nil {
		t.Error(err)
	}
	stack := transparent.NewStack()
	stack.Stack(l)
	test.BasicStackFunc(t, stack)
}

func TestFilesystemCache(t *testing.T) {
	c, err := NewCache(10, "/tmp")
	if err != nil {
		t.Error(err)
	}
	test.BasicCacheFunc(t, c)
}
