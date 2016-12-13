package filesystem

import (
	"io/ioutil"
	"os"
	"reflect"

	"github.com/juntaki/transparent"
	"github.com/juntaki/transparent/simple"
	"github.com/pkg/errors"
)

// simpleStorage store file at directory, filename is key
type simpleStorage struct {
	directory string
}

// NewSimpleStorage returns SimpleStorage
// SimpleStorage only accepts string key and []byte value.
func NewSimpleStorage(directory string) transparent.Storage {
	return &simpleStorage{
		directory: directory + "/",
	}
}

// NewStorage returns Storage
func NewStorage(directory string) transparent.Storage {
	return &simple.StorageWrapper{
		Storage: &simpleStorage{
			directory: directory + "/",
		}}
}

// Get is file read
func (f *simpleStorage) Get(k interface{}) (interface{}, error) {
	filename, err := f.validateKey(k)
	if err != nil {
		return nil, err
	}
	data, cause := ioutil.ReadFile(f.directory + filename)
	if cause != nil {
		if os.IsNotExist(cause) {
			return nil, &transparent.KeyNotFoundError{Key: filename}
		}
		return nil, errors.Wrapf(cause, "failed to read file. filename = %s", filename)
	}
	return data, nil
}

// Add is file write
func (f *simpleStorage) Add(k interface{}, v interface{}) error {
	filename, err := f.validateKey(k)
	if err != nil {
		return err
	}
	data, err := f.validateValue(v)
	if err != nil {
		return err
	}
	cause := ioutil.WriteFile(f.directory+filename, data, 0600)
	if cause != nil {
		return errors.Wrapf(cause, "failed to write file. filename = %s", filename)
	}
	return nil
}

// Remove is file unlink
func (f *simpleStorage) Remove(k interface{}) error {
	filename, err := f.validateKey(k)
	if err != nil {
		return err
	}
	cause := os.Remove(f.directory + filename)
	if cause != nil {
		return errors.Wrapf(cause, "failed to remove file. filename = %s", filename)
	}
	return nil
}

func (f *simpleStorage) validateKey(k interface{}) (string, error) {
	key, ok := k.(string)
	if !ok {
		return "", &simple.StorageInvalidKeyError{
			Valid:   reflect.TypeOf((string)("")),
			Invalid: reflect.TypeOf(k),
		}
	}
	return key, nil
}

func (f *simpleStorage) validateValue(v interface{}) ([]byte, error) {
	value, ok := v.([]byte)
	if !ok {
		return []byte{}, &simple.StorageInvalidValueError{
			Valid:   reflect.TypeOf(([]byte)("")),
			Invalid: reflect.TypeOf(v),
		}
	}
	return value, nil
}
