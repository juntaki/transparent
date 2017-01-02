package s3

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/juntaki/transparent"
	"github.com/juntaki/transparent/simple"
	"github.com/pkg/errors"
)

type BareKey struct {
	Key       string
	Bucket    string
	VersionId string
}

// Bare is key and value struct for BareStorage
type Bare struct {
	Value             map[string]interface{}
	getObjectInput    *s3.GetObjectInput
	getObjectOutput   *s3.GetObjectOutput
	putObjectInput    *s3.PutObjectInput
	deleteObjectInput *s3.DeleteObjectInput
}

// NewBare returns Bare
func NewBare() *Bare {
	return &Bare{
		Value:             map[string]interface{}{},
		getObjectInput:    &s3.GetObjectInput{},
		getObjectOutput:   &s3.GetObjectOutput{},
		putObjectInput:    &s3.PutObjectInput{},
		deleteObjectInput: &s3.DeleteObjectInput{},
	}
}

func (sb *Bare) merge(key *BareKey) {
	sb.Value["Key"] = aws.String(key.Key)
	sb.Value["Bucket"] = aws.String(key.Bucket)
	if key.VersionId != "" {
		sb.Value["VersionId"] = aws.String(key.VersionId)
	}
}

// Set Value to s3 Objects
func (sb *Bare) set() error {
	objects := []interface{}{
		sb.getObjectInput,
		sb.putObjectInput,
		sb.deleteObjectInput,
	}
	for _, o := range objects {
		ov := reflect.ValueOf(o).Elem()
		for k, v := range sb.Value {
			if k == "Body" {
				sb.putObjectInput.Body = bytes.NewReader(v.([]byte))
				continue
			}
			field := ov.FieldByName(k)
			if field.IsValid() {
				if reflect.TypeOf(v) == field.Type() {
					field.Set(reflect.ValueOf(v))
				} else {
					return fmt.Errorf("type is not matched %s %s", reflect.TypeOf(v), field.Type())
				}
			}
		}
	}
	return nil
}

// Get Value from s3 object
func (sb *Bare) get(object interface{}) {
	giv := reflect.ValueOf(object).Elem()
	for i := 0; i < giv.NumField(); i++ {
		key := giv.Type().Field(i).Name
		mapValue := giv.Field(i)
		if mapValue.CanInterface() {
			if key == "Body" {
				if mapValue.Interface() != nil {
					sb.Value[key], _ = ioutil.ReadAll(mapValue.Interface().(io.Reader))
				}
			} else {
				sb.Value[key] = mapValue.Interface()
			}
		}
	}
}

type bareStorage struct {
	svc s3iface.S3API
}

// NewBareStorage returns BareStorage
func NewBareStorage(svc s3iface.S3API) transparent.BackendStorage {
	return &bareStorage{
		svc: svc,
	}
}
func (b *bareStorage) Get(key interface{}) (value interface{}, err error) {
	bkey, err := b.validateBareKey(key)
	if err != nil {
		return nil, err
	}

	bare := NewBare()
	bare.merge(bkey)

	err = bare.set()
	if err != nil {
		return nil, err
	}

	var cause error
	bare.getObjectOutput, cause = b.svc.GetObject(bare.getObjectInput)
	if cause != nil {
		if aerr, ok := cause.(awserr.Error); ok {
			if aerr.Code() == "NoSuchKey" {
				return nil, &transparent.KeyNotFoundError{Key: key}
			}
		}
		return nil, errors.Wrapf(cause, "GetObject failed. key = %b", bare.Value)
	}
	bare.get(bare.getObjectOutput)

	return interface{}(bare), nil
}

func (b *bareStorage) Add(key interface{}, value interface{}) error {
	bkey, err := b.validateBareKey(key)
	if err != nil {
		return err
	}
	bvalue, err := b.validateBare(value)
	if err != nil {
		return err
	}

	bvalue.merge(bkey)

	err = bvalue.set()
	if err != nil {
		return err
	}

	var cause error
	_, cause = b.svc.PutObject(bvalue.putObjectInput)
	if cause != nil {
		return errors.Wrapf(cause, "PutObject failed. key = %s", bvalue.Value)
	}
	return nil
}

func (b *bareStorage) Remove(key interface{}) error {
	bkey, err := b.validateBareKey(key)
	if err != nil {
		return err
	}
	bare := NewBare()
	bare.merge(bkey)

	err = bare.set()
	if err != nil {
		return err
	}

	var cause error
	_, cause = b.svc.DeleteObject(bare.deleteObjectInput)
	if cause != nil {
		return errors.Wrapf(cause, "DeleteObject failed. key = %b", bare.Value)
	}

	return nil
}

func (b *bareStorage) validateBareKey(v interface{}) (*BareKey, error) {
	value, ok := v.(BareKey)
	if !ok {
		return nil, &simple.StorageInvalidValueError{
			Valid:   reflect.TypeOf(([]byte)("")),
			Invalid: reflect.TypeOf(v),
		}
	}
	return &value, nil
}

func (b *bareStorage) validateBare(v interface{}) (*Bare, error) {
	value, ok := v.(*Bare)
	if !ok {
		return nil, &simple.StorageInvalidValueError{
			Valid:   reflect.TypeOf(([]byte)("")),
			Invalid: reflect.TypeOf(v),
		}
	}
	return value, nil
}
