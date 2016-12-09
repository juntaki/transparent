package transparent

import (
	"errors"
	"testing"
)

func basicStorageCommand(t *testing.T, storage Storage) {
	err := storage.Add("test", []byte("value"))
	if err != nil {
		t.Error(err)
	}
	value, err := storage.Get("test")
	if err != nil || string(value.([]byte)) != "value" {
		t.Error(err)
		t.Error(value)
	}
	storage.Remove("test")
	value2, err := storage.Get("test")
	if err == nil {
		t.Error(err)
		t.Error(value2)
	}
}

func TestCustomStorage(t *testing.T) {
	test := make(map[interface{}]interface{})

	getFunc := func(k interface{}) (interface{}, error) {
		value, ok := test[k]
		if !ok {
			return nil, errors.New("value not found")
		}
		return value, nil
	}
	addFunc := func(k interface{}, v interface{}) error {
		test[k] = v
		return nil
	}
	removeFunc := func(k interface{}) error {
		delete(test, k)
		return nil
	}

	storage, err := NewCustomStorage(getFunc, addFunc, removeFunc)
	if err != nil {
		t.Error(err)
	}
	basicStorageCommand(t, storage)
}

func TestFilesystemStorage(t *testing.T) {
	storage, err := NewFilesystemStorage("/tmp")
	if err != nil {
		t.Error(err)
	}
	basicStorageCommand(t, storage)
}

func TestS3Storage(t *testing.T) {
	svc, err := newMockS3Client()
	if err != nil {
		t.Error(err)
	}
	storage, err := NewS3Storage("bucket", svc)
	if err != nil {
		t.Error(err)
	}
	basicStorageCommand(t, storage)
}
