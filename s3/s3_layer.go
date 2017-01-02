package s3

import (
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/juntaki/transparent"
)

// NewS3Cache returns S3Cache
func NewCache(bufferSize int, bucket string, svc s3iface.S3API) (transparent.Layer, error) {
	s3 := NewSimpleStorage(bucket, svc)
	layer, err := transparent.NewLayerCache(bufferSize, s3)
	if err != nil {
		return nil, err
	}
	return layer, nil
}

// NewS3Source returns S3Source
func NewSource(bucket string, svc s3iface.S3API) (transparent.Layer, error) {
	s3 := NewSimpleStorage(bucket, svc)
	layer, err := transparent.NewLayerSource(s3)
	if err != nil {
		return nil, err
	}
	return layer, nil
}

func NewBareSource(svc s3iface.S3API) (transparent.Layer, error) {
	s3 := NewBareStorage(svc)
	layer, err := transparent.NewLayerSource(s3)
	if err != nil {
		return nil, err
	}
	return layer, nil
}
