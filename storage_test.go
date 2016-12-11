package transparent

import (
	"reflect"
	"testing"
)

// Get, Add and Remove
func basicStorageFunc(t *testing.T, storage Storage) {
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
	storageErr, ok := err.(*StorageKeyNotFoundError)
	if ok {
		if storageErr.Key != "test" {
			t.Fatal("key is different", storageErr.Key)
		}
	} else {
		t.Fatal(err, value2)
	}
}

func simpleStorageFunc(t *testing.T, storage Storage) {
	var err error
	err = storage.Add(0, []byte("value"))
	if typeErr, ok := err.(*SimpleStorageInvalidKeyError); !ok {
		t.Fatal(err)
	} else if typeErr.invalid != reflect.TypeOf(0) {
		t.Fatal(typeErr)
	}
	err = storage.Add("test", 0)
	if typeErr, ok := err.(*SimpleStorageInvalidValueError); !ok {
		t.Fatal(err)
	} else if typeErr.invalid != reflect.TypeOf(0) {
		t.Fatal(typeErr)
	}

	_, err = storage.Get(0)
	if typeErr, ok := err.(*SimpleStorageInvalidKeyError); !ok {
		t.Fatal(err)
	} else if typeErr.invalid != reflect.TypeOf(0) {
		t.Fatal(typeErr)
	}

	err = storage.Remove(0)
	if typeErr, ok := err.(*SimpleStorageInvalidKeyError); !ok {
		t.Fatal(err)
	} else if typeErr.invalid != reflect.TypeOf(0) {
		t.Fatal(typeErr)
	}
}

func TestStorage(t *testing.T) {
	var err error

	// Dummy
	ds, err := NewDummyStorage(0)
	if err != nil {
		t.Fatal(err)
	}
	basicStorageFunc(t, ds)

	// Custom
	ds, err = NewDummyStorage(0)
	if err != nil {
		t.Fatal(err)
	}
	cs, err := NewCustomStorage(ds.Get, ds.Add, ds.Remove)
	if err != nil {
		t.Fatal(err)
	}
	basicStorageFunc(t, cs)

	fss, err := NewFilesystemSimpleStorage("/tmp")
	if err != nil {
		t.Fatal(err)
	}
	basicStorageFunc(t, fss)
	simpleStorageFunc(t, fss)

	// Filesystem
	fs, err := NewFilesystemStorage("/tmp")
	if err != nil {
		t.Fatal(err)
	}
	basicStorageFunc(t, fs)

	// S3
	svc, err := newMockS3Client()
	if err != nil {
		t.Fatal(err)
	}
	sss, err := NewS3SimpleStorage("bucket", svc)
	if err != nil {
		t.Fatal(err)
	}
	basicStorageFunc(t, sss)
	simpleStorageFunc(t, sss)

	svc, err = newMockS3Client()
	if err != nil {
		t.Fatal(err)
	}
	ss, err := NewS3Storage("bucket", svc)
	if err != nil {
		t.Fatal(err)
	}
	basicStorageFunc(t, ss)
}
