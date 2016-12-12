package transparent_test

import (
	"reflect"
	"testing"

	"github.com/juntaki/transparent"
	"github.com/juntaki/transparent/dummy"
	"github.com/juntaki/transparent/simple"
)

func SimpleStorageFunc(t *testing.T, storage transparent.Storage) {
	var err error
	err = storage.Add(0, []byte("value"))
	if typeErr, ok := err.(*simple.StorageInvalidKeyError); !ok {
		t.Fatal(err)
	} else if typeErr.Invalid != reflect.TypeOf(0) {
		t.Fatal(typeErr)
	}
	err = storage.Add("test", 0)
	if typeErr, ok := err.(*simple.StorageInvalidValueError); !ok {
		t.Fatal(err)
	} else if typeErr.Invalid != reflect.TypeOf(0) {
		t.Fatal(typeErr)
	}

	_, err = storage.Get(0)
	if typeErr, ok := err.(*simple.StorageInvalidKeyError); !ok {
		t.Fatal(err)
	} else if typeErr.Invalid != reflect.TypeOf(0) {
		t.Fatal(typeErr)
	}

	err = storage.Remove(0)
	if typeErr, ok := err.(*simple.StorageInvalidKeyError); !ok {
		t.Fatal(err)
	} else if typeErr.Invalid != reflect.TypeOf(0) {
		t.Fatal(typeErr)
	}
}

// Get, Add and Remove
func BasicStorageFunc(t *testing.T, storage transparent.Storage) {
	// Add and Get
	err := storage.Add("test", []byte("value"))
	if err != nil {
		t.Fatal(err)
	}
	value, err := storage.Get("test")
	if err != nil || string(value.([]byte)) != "value" {
		t.Fatal(err, value)
	}

	// Remove and Get
	storage.Remove("test")
	value2, err := storage.Get("test")
	storageErr, ok := err.(*transparent.StorageKeyNotFoundError)
	if ok {
		if storageErr.Key != "test" {
			t.Fatal("key is different", storageErr.Key)
		}
	} else {
		t.Fatal(err, value2)
	}
}

func BasicStackFunc(t *testing.T, l *transparent.Stack) {
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

func BasicCacheFunc(t *testing.T, c *transparent.Cache) {
	s, err := dummy.NewSource(0)
	if err != nil {
		t.Error(err)
	}

	stack := transparent.NewStack()
	stack.Stack(s)
	stack.Stack(c)
	stack.Start()
	BasicStackFunc(t, stack)
	stack.Stop()
}
