package s3

import (
	"reflect"

	"github.com/juntaki/transparent"
	"github.com/juntaki/transparent/simple"
	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// NewS3Storage returns S3Storage
func NewStorage(bucket string, svc s3iface.S3API) transparent.BackendStorage {
	s := NewSimpleStorage(bucket, svc)
	return &simple.StorageWrapper{
		BackendStorage: s,
	}
}

// s3SimpleStorage store file to Amazon S3 as object
type simpleStorage struct {
	bare   transparent.BackendStorage
	svc    s3iface.S3API
	bucket string
}

// NewS3SimpleStorage returns s3SimpleStorage
func NewSimpleStorage(bucket string, svc s3iface.S3API) transparent.BackendStorage {
	return &simpleStorage{
		bare:   NewBareStorage(svc),
		svc:    svc,
		bucket: bucket,
	}
}

// Get is get request
func (s *simpleStorage) Get(k interface{}) (interface{}, error) {
	key, err := s.validateKey(k)
	if err != nil {
		return nil, err
	}

	bk := BareKey{
		Key:    key,
		Bucket: s.bucket,
	}

	br, err := s.bare.Get(bk)
	if err != nil {
		if _, ok := err.(*transparent.KeyNotFoundError); ok {
			return nil, &transparent.KeyNotFoundError{Key: key}
		}
		return nil, err
	}

	body := br.(*Bare).Value["Body"]
	return body, nil
}

// Add is set put request
func (s *simpleStorage) Add(k interface{}, v interface{}) error {
	key, err := s.validateKey(k)
	if err != nil {
		return err
	}
	body, err := s.validateValue(v)
	if err != nil {
		return err
	}

	bk := BareKey{
		Key:    key,
		Bucket: s.bucket,
	}

	bv := NewBare()
	bv.Value["Body"] = body

	return s.bare.Add(bk, bv)
}

// Remove is delete request
func (s *simpleStorage) Remove(k interface{}) error {
	key, err := s.validateKey(k)
	if err != nil {
		return err
	}
	params := &s3.DeleteObjectInput{
		Key:    aws.String(key),
		Bucket: aws.String(s.bucket),
	}
	_, cause := s.svc.DeleteObject(params)
	if cause != nil {
		return errors.Wrapf(cause, "DeleteObject failed. key = %s", key)
	}
	return nil
}

func (s *simpleStorage) validateKey(k interface{}) (string, error) {
	key, ok := k.(string)
	if !ok {
		return "", &simple.StorageInvalidKeyError{
			Valid:   reflect.TypeOf((string)("")),
			Invalid: reflect.TypeOf(k),
		}
	}
	return key, nil
}

func (s *simpleStorage) validateValue(v interface{}) ([]byte, error) {
	value, ok := v.([]byte)
	if !ok {
		return []byte{}, &simple.StorageInvalidValueError{
			Valid:   reflect.TypeOf(([]byte)("")),
			Invalid: reflect.TypeOf(v),
		}
	}
	return value, nil
}
