package s3

import (
	"bytes"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/juntaki/transparent"
	"github.com/juntaki/transparent/test"
)

type mockS3Client struct {
	s3iface.S3API
	d transparent.BackendStorage
}

func newMockS3Client() (*mockS3Client, error) {
	test := test.NewStorage(0)
	return &mockS3Client{d: test}, nil
}

func (m *mockS3Client) GetObject(i *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	if *i.Bucket != "bucket" {
		return nil, errors.New("bucket name invalid")
	}
	value, err := m.d.Get(*i.Key)
	if err != nil {
		aerr := awserr.New("NoSuchKey", "NoSuchKeyDummy", err)
		return nil, aerr
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

func TestStorage(t *testing.T) {
	var err error

	// S3
	svc, err := newMockS3Client()
	if err != nil {
		t.Fatal(err)
	}
	sss, err := NewSimpleStorage("bucket", svc)
	if err != nil {
		t.Fatal(err)
	}
	test.BasicStorageFunc(t, sss)
	test.SimpleStorageFunc(t, sss)

	svc, err = newMockS3Client()
	if err != nil {
		t.Fatal(err)
	}
	ss, err := NewStorage("bucket", svc)
	if err != nil {
		t.Fatal(err)
	}
	test.BasicStorageFunc(t, ss)

	svc, err = newMockS3Client()
	if err != nil {
		t.Error(err)
	}
	l, err := NewSource("bucket", svc)
	if err != nil {
		t.Error(err)
	}

	stack := transparent.NewStack()
	stack.Stack(l)
	test.BasicStackFunc(t, stack)

	svc, err = newMockS3Client()
	if err != nil {
		t.Error(err)
	}
	c, err := NewCache(10, "bucket", svc)
	if err != nil {
		t.Error(err)
	}
	test.BasicCacheFunc(t, c)
}
