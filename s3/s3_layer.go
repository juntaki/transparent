package s3

import (
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/juntaki/transparent"
)

// NewS3Cache returns S3Cache
func NewCache(bufferSize int, bucket string, svc s3iface.S3API) (*transparent.LayerCache, error) {
	s3, err := NewSimpleStorage(bucket, svc)
	if err != nil {
		return nil, err
	}
	layer, err := transparent.NewLayerCache(bufferSize, s3)
	if err != nil {
		return nil, err
	}
	return layer, nil
}

// NewS3Source returns S3Source
func NewSource(bucket string, svc s3iface.S3API) (*transparent.LayerSource, error) {
	s3, err := NewSimpleStorage(bucket, svc)
	if err != nil {
		return nil, err
	}
	layer, err := transparent.NewLayerSource(s3)
	if err != nil {
		return nil, err
	}
	return layer, nil
}
