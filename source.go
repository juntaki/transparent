package transparent

import (
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// Source provides operation of TransparentSource
type Source struct {
	Storage Storage
	upper   Layer
}

// NewSource returns Source
func NewSource(storage Storage) (*Source, error) {
	if storage == nil {
		return nil, errors.New("empty storage")
	}
	return &Source{Storage: storage}, nil
}

// Set set new value to storage.
func (s *Source) Set(key interface{}, value interface{}) (err error) {
	err = s.Storage.Add(key, value)
	if err != nil {
		return err
	}
	return nil
}

// Get value from storage
func (s *Source) Get(key interface{}) (value interface{}, err error) {
	return s.Storage.Get(key)
}

// Remove value
func (s *Source) Remove(key interface{}) (err error) {
	return s.Storage.Remove(key)
}

// Sync do nothing
func (s *Source) Sync() error {
	return nil
}

func (s *Source) setUpper(upper Layer) {
	s.upper = upper
}

func (s *Source) setLower(lower Layer) {
	panic("don't set lower layer")
}

// NewDummySource returns dummyStorage layer
func NewDummySource(wait time.Duration) (*Source, error) {
	dummy, err := NewDummyStorage(wait)
	if err != nil {
		return nil, err
	}
	layer, _ := NewSource(dummy)
	return layer, nil
}

// NewS3Source returns S3Source
func NewS3Source(bucket string, svc s3iface.S3API) (*Source, error) {
	s3, err := NewS3SimpleStorage(bucket, svc)
	if err != nil {
		return nil, err
	}
	layer, err := NewSource(s3)
	if err != nil {
		return nil, err
	}
	return layer, nil
}

// NewFilesystemSource returns FilesystemSource
func NewFilesystemSource(directory string) (*Source, error) {
	filesystem, err := NewFilesystemSimpleStorage(directory)
	if err != nil {
		return nil, err
	}
	layer, err := NewSource(filesystem)
	if err != nil {
		return nil, err
	}
	return layer, nil
}
