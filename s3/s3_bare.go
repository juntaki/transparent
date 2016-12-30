package s3

import (
	"reflect"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/juntaki/transparent"
	"github.com/juntaki/transparent/simple"
	"github.com/pkg/errors"
)

type BareKey struct {
	Bucket *string
	Key    *string
}

type BareValue struct {
	getObjectInput  *s3.GetObjectInput
	getObjectOutput *s3.GetObjectOutput
	putObjectInput  *s3.PutObjectInput
	//putObjectOutput    *s3.PutObjectOutput
	deleteObjectInput *s3.DeleteObjectInput
	//deleteObjectOutput *s3.DeleteObjectOutput
	Key           *BareKey
	Value         map[string]interface{}
	lastOperatoin bool
}

type bare struct {
	*BareValue
	Key *BareKey
}

func NewBareValue() *BareValue {
	return &BareValue{
		Value:           map[string]interface{}{},
		getObjectInput:  &s3.GetObjectInput{},
		getObjectOutput: &s3.GetObjectOutput{},
		putObjectInput:  &s3.PutObjectInput{},
		//putObjectOutput:    &s3.PutObjectOutput{},
		deleteObjectInput: &s3.DeleteObjectInput{},
		//deleteObjectOutput: &s3.DeleteObjectOutput{},
	}
}

func newBare(key *BareKey, value *BareValue) *bare {
	return &bare{
		BareValue: value,
		Key:       key,
	}
}

// Set Value to s3 Objects
func (sb *bare) Set() error {
	objects := []interface{}{
		sb.getObjectOutput,
		sb.getObjectInput,
		//sb.putObjectOutput,
		sb.putObjectInput,
		//sb.deleteObjectOutput,
		sb.deleteObjectInput,
	}
	for _, o := range objects {
		ov := reflect.ValueOf(o).Elem()
		for k, v := range sb.Value {
			field := ov.FieldByName(k)
			if field.IsValid() {
				if reflect.TypeOf(v) == field.Type() {
					field.Set(reflect.ValueOf(v))
				} else {
					return errors.New("type is not matched")
				}
			}
		}
	}
	sb.getObjectInput.Key = sb.Key.Key
	sb.getObjectInput.Bucket = sb.Key.Bucket
	sb.putObjectInput.Key = sb.Key.Key
	sb.putObjectInput.Bucket = sb.Key.Bucket
	sb.deleteObjectInput.Key = sb.Key.Key
	sb.deleteObjectInput.Bucket = sb.Key.Bucket
	return nil
}

// Get Value from s3 object
func (sb *bare) Get(object interface{}) {
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

func NewbareStorage(svc s3iface.S3API) transparent.BackendStorage {
	return &bareStorage{
		svc: svc,
	}
}
func (b *bareStorage) Get(key interface{}) (value interface{}, err error) {
	bkey, err := b.validateKey(key)
	if err != nil {
		return nil, err
	}
	bvalue := NewBareValue()
	bare := newBare(&bkey, bvalue)

	var cause error
	bare.getObjectOutput, cause = b.svc.GetObject(bare.getObjectInput)
	if cause != nil {
		if aerr, ok := cause.(awserr.Error); ok {
			if aerr.Code() == "NoSuchKey" {
				return nil, &transparent.KeyNotFoundError{Key: key}
			}
		}
		return nil, errors.Wrapf(cause, "GetObject failed. key = %b", bkey)
	}

	bare.Get(bare.getObjectInput)
	bare.Get(bare.getObjectOutput)

	return interface{}(bare.BareValue), nil
}

func (b *bareStorage) Add(key interface{}, value interface{}) error {
	bkey, err := b.validateKey(key)
	if err != nil {
		return err
	}
	bvalue, err := b.validateValue(value)
	if err != nil {
		return err
	}

	bare := newBare(&bkey, bvalue)
	err = bare.Set()
	if err != nil {
		return err
	}

	var cause error
	_, cause = b.svc.PutObject(bare.putObjectInput)
	if cause != nil {
		return errors.Wrapf(cause, "PutObject failed. key = %s", bkey)
	}

	return nil
}

func (b *bareStorage) Remove(key interface{}) error {
	bkey, err := b.validateKey(key)
	if err != nil {
		return err
	}
	bvalue := NewBareValue()
	bare := newBare(&bkey, bvalue)
	err = bare.Set()
	if err != nil {
		return err
	}

	var cause error
	_, cause = b.svc.DeleteObject(bare.deleteObjectInput)
	if cause != nil {
		return errors.Wrapf(cause, "DeleteObject failed. key = %b", bkey)
	}

	return nil
}
func (b *bareStorage) validateKey(k interface{}) (BareKey, error) {
	key, ok := k.(BareKey)
	if !ok {
		return BareKey{}, &simple.StorageInvalidKeyError{
			Valid:   reflect.TypeOf((string)("")),
			Invalid: reflect.TypeOf(k),
		}
	}
	return key, nil
}

func (b *bareStorage) validateValue(v interface{}) (*BareValue, error) {
	value, ok := v.(*BareValue)
	if !ok {
		return nil, &simple.StorageInvalidValueError{
			Valid:   reflect.TypeOf(([]byte)("")),
			Invalid: reflect.TypeOf(v),
		}
	}
	return value, nil
}
