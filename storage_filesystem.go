package transparent

import (
	"io/ioutil"
	"os"
	"reflect"

	"github.com/pkg/errors"
)

// FilesystemStorage store file at directory, filename is key
type FilesystemStorage struct {
	directory string
}

// NewFilesystemStorage returns FilesystemStorage
func NewFilesystemStorage(directory string) (*FilesystemStorage, error) {
	return &FilesystemStorage{
		directory: directory + "/",
	}, nil
}

// Get is file read
func (f *FilesystemStorage) Get(k interface{}) (interface{}, error) {
	filename, err := f.validateKey(k)
	if err != nil {
		return nil, err
	}
	data, cause := ioutil.ReadFile(f.directory + filename)
	if cause != nil {
		if os.IsNotExist(cause) {
			return nil, &StorageKeyNotFoundError{Key: filename}
		}
		return nil, errors.Wrapf(cause, "failed to read file. filename = %s", filename)
	}
	return data, nil
}

// Add is file write
func (f *FilesystemStorage) Add(k interface{}, v interface{}) error {
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
func (f *FilesystemStorage) Remove(k interface{}) error {
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

func (f *FilesystemStorage) validateKey(k interface{}) (string, error) {
	key, ok := k.(string)
	if !ok {
		return "", &StorageInvalidKeyError{
			valid:   reflect.TypeOf((string)("")),
			invalid: reflect.TypeOf(k),
		}
	}
	return key, nil
}

func (f *FilesystemStorage) validateValue(v interface{}) ([]byte, error) {
	value, ok := v.([]byte)
	if !ok {
		return []byte{}, &StorageInvalidValueError{
			valid:   reflect.TypeOf(([]byte)("")),
			invalid: reflect.TypeOf(v),
		}
	}
	return value, nil
}
