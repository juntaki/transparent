package test

import (
	"reflect"
	"testing"

	"github.com/juntaki/transparent"
	"github.com/juntaki/transparent/simple"
)

// SimpleStorageFunc is Error message test for simple Storage
func SimpleStorageFunc(t *testing.T, storage transparent.BackendStorage) {
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

// BasicStorageFunc is Get, Add and Remove
func BasicStorageFunc(t *testing.T, storage transparent.BackendStorage) {
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

// BasicStackFunc is Get Remove and Sync
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

// BasicCacheFunc is test for transparent.Cache
func BasicCacheFunc(t *testing.T, c *transparent.LayerCache) {
	s := NewSource(0)
	stack := transparent.NewStack()
	stack.Stack(s)
	stack.Stack(c)
	stack.Start()
	BasicStackFunc(t, stack)
	stack.Stop()
}

// BasicSourceFunc is test for transparent.Source
func BasicSourceFunc(t *testing.T, s *transparent.LayerSource) {
	stack := transparent.NewStack()
	stack.Stack(s)
	stack.Start()
	BasicStackFunc(t, stack)
	stack.Stop()
}

// BasicTransmitterFunc is test for transparent.Transmitter
func BasicTransmitterFunc(t *testing.T, s *transparent.LayerTransmitter) {
	stack := transparent.NewStack()
	stack.Stack(s)
	stack.Start()
	BasicStackFunc(t, stack)
	stack.Stop()
}

// BasicConsensusFunc is test for transparent.Consensus
func BasicConsensusFunc(t *testing.T, a1, a2 *transparent.LayerConsensus) {
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
