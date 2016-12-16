package s3

import (
	"bytes"
	"io/ioutil"
	"reflect"

	"github.com/juntaki/transparent"
	"github.com/juntaki/transparent/simple"
	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// NewS3Storage returns S3Storage
func NewStorage(bucket string, svc s3iface.S3API) (transparent.BackendStorage, error) {
	return &simple.StorageWrapper{
		BackendStorage: &simpleStorage{
			svc:    svc,
			bucket: aws.String(bucket),
		}}, nil
}

// s3SimpleStorage store file to Amazon S3 as object
type simpleStorage struct {
	svc    s3iface.S3API
	bucket *string
}

// NewS3SimpleStorage returns s3SimpleStorage
func NewSimpleStorage(bucket string, svc s3iface.S3API) (transparent.BackendStorage, error) {
	return &simpleStorage{
		svc:    svc,
		bucket: aws.String(bucket),
	}, nil
}

// Get is get request
func (s *simpleStorage) Get(k interface{}) (interface{}, error) {
	key, err := s.validateKey(k)
	if err != nil {
		return nil, err
	}
	paramsGet := &s3.GetObjectInput{
		Bucket: s.bucket,
		Key:    aws.String(key),
	}
	respGet, cause := s.svc.GetObject(paramsGet)

	if cause != nil {
		if aerr, ok := cause.(awserr.Error); ok {
			if aerr.Code() == "NoSuchKey" {
				return nil, &transparent.KeyNotFoundError{Key: key}
			}
		}
		return nil, errors.Wrapf(cause, "GetObject failed. key = %s", key)
	}
	body, cause := ioutil.ReadAll(respGet.Body)
	if cause != nil {
		return nil, errors.Wrapf(cause, "failed to read response body. key = %s", key)
	}
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

	params := &s3.PutObjectInput{
		Bucket: s.bucket,
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
	}
	_, cause := s.svc.PutObject(params)
	if cause != nil {
		return errors.Wrapf(cause, "PutObject failed. key = %s", key)
	}
	return nil
}

// Remove is delete request
func (s *simpleStorage) Remove(k interface{}) error {
	key, err := s.validateKey(k)
	if err != nil {
		return err
	}
	params := &s3.DeleteObjectInput{
		Bucket: s.bucket,
		Key:    aws.String(key),
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
