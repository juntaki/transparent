package transparent

import (
	"errors"
	"io/ioutil"
	"os"
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
	filename, ok := k.(string)
	if !ok {
		return nil, errors.New("key must be string")
	}
	data, err := ioutil.ReadFile(f.directory + filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Add is file write
func (f *FilesystemStorage) Add(k interface{}, v interface{}) error {
	filename, ok := k.(string)
	if !ok {
		return errors.New("key must be string")
	}
	data, ok := v.([]byte)
	if !ok {
		return errors.New("value must be []byte")
	}
	err := ioutil.WriteFile(f.directory+filename, data, 0600)
	return err
}

// Remove is file unlink
func (f *FilesystemStorage) Remove(k interface{}) error {
	filename, ok := k.(string)
	if !ok {
		return errors.New("key must be string")
	}
	return os.Remove(f.directory + filename)
}
