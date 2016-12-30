package s3

import (
	"fmt"
	"io"
	"reflect"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/juntaki/transparent"
	"github.com/juntaki/transparent/simple"
	"github.com/pkg/errors"
)

// Bare is key and value struct for BareStorage
type Bare struct {
	Value             map[string]interface{}
	getObjectInput    *s3.GetObjectInput
	getObjectOutput   *s3.GetObjectOutput
	putObjectInput    *s3.PutObjectInput
	deleteObjectInput *s3.DeleteObjectInput
}

// NewBare returns Bare
func NewBare() Bare {
	return Bare{
		Value:             map[string]interface{}{},
		getObjectInput:    &s3.GetObjectInput{},
		getObjectOutput:   &s3.GetObjectOutput{},
		putObjectInput:    &s3.PutObjectInput{},
		deleteObjectInput: &s3.DeleteObjectInput{},
	}
}

func (sb *Bare) merge(target *Bare) error {
	for k, v := range target.Value {
		if _, ok := sb.Value[k]; ok {
			return errors.New("Duplicate value")
		}
		sb.Value[k] = v
	}
	return nil
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
				sb.putObjectInput.Body = v.(io.ReadSeeker)
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
			sb.Value[key] = mapValue.Interface()
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
	bare, err := b.validateBare(key)
	if err != nil {
		return nil, err
	}
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
	bkey, err := b.validateBare(key)
	if err != nil {
		return err
	}
	bvalue, err := b.validateBare(value)
	if err != nil {
		return err
	}

	err = bkey.merge(bvalue)
	if err != nil {
		return err
	}

	err = bkey.set()
	if err != nil {
		return err
	}

	var cause error
	_, cause = b.svc.PutObject(bkey.putObjectInput)
	if cause != nil {
		return errors.Wrapf(cause, "PutObject failed. key = %s", bkey.Value)
	}
	return nil
}

func (b *bareStorage) Remove(key interface{}) error {
	bare, err := b.validateBare(key)
	if err != nil {
		return err
	}
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

func (b *bareStorage) validateBare(v interface{}) (*Bare, error) {
	value, ok := v.(Bare)
	if !ok {
		return nil, &simple.StorageInvalidValueError{
			Valid:   reflect.TypeOf(([]byte)("")),
			Invalid: reflect.TypeOf(v),
		}
	}
	return &value, nil
}
