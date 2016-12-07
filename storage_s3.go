package transparent

import (
	"bytes"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Storage store file at directory, filename is key
type S3Storage struct {
	svc    *s3.S3
	bucket *string
}

// NewS3Storage returns S3Storage
func NewS3Storage(bucket string, cfgs ...*aws.Config) (*S3Storage, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	svc := s3.New(sess, cfgs...)
	return &S3Storage{
		svc:    svc,
		bucket: aws.String(bucket),
	}, nil
}

// Get is get request
func (s *S3Storage) Get(k interface{}) (interface{}, error) {
	key, ok := k.(string)
	if !ok {
		return nil, errors.New("key must be string")
	}
	paramsGet := &s3.GetObjectInput{
		Bucket: s.bucket,
		Key:    aws.String(key),
	}
	respGet, err := s.svc.GetObject(paramsGet)
	if err != nil {
		return nil, err
	}
	return respGet.Body, nil
}

// Add is set put request
func (s *S3Storage) Add(k interface{}, v interface{}) error {
	key, ok := k.(string)
	if !ok {
		return errors.New("key must be string")
	}
	body, ok := v.([]byte)
	if !ok {
		return errors.New("value must be []byte")
	}

	params := &s3.PutObjectInput{
		Bucket: s.bucket,
		Key:    aws.String(key),
		Body:   bytes.NewReader(body),
	}
	_, err := s.svc.PutObject(params)
	if err != nil {
		return err
	}
	return nil
}

// Remove is delete request
func (s *S3Storage) Remove(k interface{}) error {
	key, ok := k.(string)
	if !ok {
		return errors.New("key must be string")
	}
	params := &s3.DeleteObjectInput{
		Bucket: s.bucket,
		Key:    aws.String(key),
	}
	_, err := s.svc.DeleteObject(params)
	if err != nil {
		return err
	}
	return nil
}
