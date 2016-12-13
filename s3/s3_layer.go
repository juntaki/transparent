package s3

import (
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/juntaki/transparent"
)

// NewS3Cache returns S3Cache
func NewCache(bufferSize int, bucket string, svc s3iface.S3API) (*transparent.Cache, error) {
	s3, err := NewSimpleStorage(bucket, svc)
	if err != nil {
		return nil, err
	}
	layer, err := transparent.NewCache(bufferSize, s3)
	if err != nil {
		return nil, err
	}
	return layer, nil
}

// NewS3Source returns S3Source
func NewSource(bucket string, svc s3iface.S3API) (*transparent.Source, error) {
	s3, err := NewSimpleStorage(bucket, svc)
	if err != nil {
		return nil, err
	}
	layer, err := transparent.NewSource(s3)
	if err != nil {
		return nil, err
	}
	return layer, nil
}
