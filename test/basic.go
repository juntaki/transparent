package test

import (
	"reflect"
	"testing"

	"github.com/juntaki/transparent"
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
	storageErr, ok := err.(*transparent.KeyNotFoundError)
	if ok {
		if storageErr.Key != "test" {
			t.Fatal("key is different", storageErr.Key)
		}
	} else {
		t.Fatal(err, value2)
	}
}

func BasicStackFunc(t *testing.T, s *transparent.Stack) {
	err := s.Set("test", []byte("value"))
	if err != nil {
		t.Error(err)
	}

	value, err := s.Get("test")
	if err != nil || string(value.([]byte)) != "value" {
		t.Error(err)
		t.Error(value)
	}

	err = s.Remove("test")
	if err != nil {
		t.Error(err)
	}

	value, err = s.Get("test")
	if err == nil {
		t.Error(err)
		t.Error(value)
	}

	err = s.Sync()
	if err != nil {
		t.Error(err)
	}
}

func BasicCacheFunc(t *testing.T, c *transparent.Cache) {
	s := NewSource(0)
	stack := transparent.NewStack()
	stack.Stack(s)
	stack.Stack(c)
	stack.Start()
	BasicStackFunc(t, stack)
	stack.Stop()
}

func BasicSourceFunc(t *testing.T, s *transparent.Source) {
	stack := transparent.NewStack()
	stack.Stack(s)
	stack.Start()
	BasicStackFunc(t, stack)
	stack.Stop()
}

func BasicConsensusFunc(t *testing.T, a1, a2 *transparent.Consensus) {
	src1 := NewSource(0)
	src2 := NewSource(0)

	s1 := transparent.NewStack()
	s2 := transparent.NewStack()

	s1.Stack(src1)
	s2.Stack(src2)

	s1.Stack(a1)
	s2.Stack(a2)

	s1.Start()
	s2.Start()

	BasicStackFunc(t, s1)
	BasicStackFunc(t, s2)

	s1.Stop()
	s2.Stop()
}
