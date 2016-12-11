package transparent

import (
	"bytes"
	"io/ioutil"
	"reflect"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// NewS3Storage returns S3Storage
func NewS3Storage(bucket string, svc s3iface.S3API) (Storage, error) {
	return &simpleStorageWrapper{
		Storage: &s3SimpleStorage{
			svc:    svc,
			bucket: aws.String(bucket),
		}}, nil
}

// s3SimpleStorage store file to Amazon S3 as object
type s3SimpleStorage struct {
	svc    s3iface.S3API
	bucket *string
}

// NewS3SimpleStorage returns s3SimpleStorage
func NewS3SimpleStorage(bucket string, svc s3iface.S3API) (Storage, error) {
	return &s3SimpleStorage{
		svc:    svc,
		bucket: aws.String(bucket),
	}, nil
}

// Get is get request
func (s *s3SimpleStorage) Get(k interface{}) (interface{}, error) {
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
				return nil, &StorageKeyNotFoundError{Key: key}
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
func (s *s3SimpleStorage) Add(k interface{}, v interface{}) error {
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
func (s *s3SimpleStorage) Remove(k interface{}) error {
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

func (s *s3SimpleStorage) validateKey(k interface{}) (string, error) {
	key, ok := k.(string)
	if !ok {
		return "", &SimpleStorageInvalidKeyError{
			valid:   reflect.TypeOf((string)("")),
			invalid: reflect.TypeOf(k),
		}
	}
	return key, nil
}

func (s *s3SimpleStorage) validateValue(v interface{}) ([]byte, error) {
	value, ok := v.([]byte)
	if !ok {
		return []byte{}, &SimpleStorageInvalidValueError{
			valid:   reflect.TypeOf(([]byte)("")),
			invalid: reflect.TypeOf(v),
		}
	}
	return value, nil
}
