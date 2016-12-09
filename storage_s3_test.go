package transparent

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type mockS3Client struct {
	s3iface.S3API
	d *DummyStorage
}

func (m *mockS3Client) GetObject(i *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	if *i.Bucket != "bucket" {
		return nil, errors.New("bucket name invalid")
	}
	value, err := m.d.Get(*i.Key)
	if err != nil {
		return nil, err
	}
	body, ok := value.([]byte)
	if !ok {
		return nil, errors.New("value invalid")
	}
	return &s3.GetObjectOutput{Body: ioutil.NopCloser(bytes.NewReader(body))}, nil
}

func (m *mockS3Client) PutObject(i *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	if *i.Bucket != "bucket" {
		return nil, errors.New("bucket name invalid")
	}
	body, err := ioutil.ReadAll(i.Body)
	if err != nil {
		return nil, err
	}
	err = m.d.Add(*i.Key, body)
	if err != nil {
		return nil, err
	}
	return &s3.PutObjectOutput{}, nil
}

func (m *mockS3Client) DeleteObject(i *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	if *i.Bucket != "bucket" {
		return nil, errors.New("bucket name invalid")
	}
	err := m.d.Remove(*i.Key)
	if err != nil {
		return nil, err
	}
	return &s3.DeleteObjectOutput{}, nil
}

func TestS3Storage(t *testing.T) {
	dummy, err := NewDummyStorage(1)
	if err != nil {
		t.Error(err)
	}

	storage, err := NewS3Storage("bucket", &mockS3Client{d: dummy})
	if err != nil {
		t.Error(err)
	}
	basicCommand(t, storage)
}