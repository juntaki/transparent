package simple

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"reflect"

	"github.com/juntaki/transparent"
	"github.com/pkg/errors"
)

type StorageWrapper struct {
	transparent.Storage
}

// Get is file read
func (f *StorageWrapper) Get(k interface{}) (interface{}, error) {
	key, err := f.encodeKey(k)
	if err != nil {
		return nil, err
	}
	v, err := f.Storage.Get(key)
	if err != nil {
		_, ok := err.(*transparent.KeyNotFoundError)
		if ok {
			return nil, &transparent.KeyNotFoundError{Key: k}
		}
		return nil, err
	}
	data, err := f.decodeValue(v.([]byte))
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Add is file write
func (f *StorageWrapper) Add(k interface{}, v interface{}) error {
	key, err := f.encodeKey(k)
	if err != nil {
		return err
	}
	data, err := f.encodeValue(v)
	if err != nil {
		return err
	}
	err = f.Storage.Add(key, data)
	if err != nil {
		return err
	}
	return nil
}

// Remove is file unlink
func (f *StorageWrapper) Remove(k interface{}) error {
	key, err := f.encodeKey(k)
	if err != nil {
		return err
	}
	err = f.Storage.Remove(key)
	if err != nil {
		return err
	}
	return nil
}

func (f *StorageWrapper) encodeKey(k interface{}) (string, error) {
	gob.Register(k)
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	cause := encoder.Encode(k)
	if cause != nil {
		return "", errors.Wrap(cause, "failed to encode key")
	}
	key := fmt.Sprintf("%x", sha256.Sum256(buf.Bytes()))
	return key, nil
}

func (f *StorageWrapper) encodeValue(v interface{}) ([]byte, error) {
	gob.Register(v)
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	cause := encoder.Encode(&v)
	if cause != nil {
		return nil, errors.Wrap(cause, "failed to encode value")
	}
	return buf.Bytes(), nil
}

func (f *StorageWrapper) decodeValue(v []byte) (interface{}, error) {
	var d interface{}
	buf := bytes.NewBuffer(v)
	decoder := gob.NewDecoder(buf)
	cause := decoder.Decode(&d)
	if cause != nil {
		return nil, errors.Wrap(cause, "failed to decode value")
	}
	return d, nil
}

// StorageInvalidKeyError means type of key is invalid for storage
type StorageInvalidKeyError struct {
	Valid   reflect.Type
	Invalid reflect.Type
}

func (e *StorageInvalidKeyError) Error() string {
	return fmt.Sprintf("%s is not supported key in the storage, use %s", e.Invalid, e.Valid)
}

// StorageInvalidValueError means type of value is invalid for storage
type StorageInvalidValueError struct {
	Valid   reflect.Type
	Invalid reflect.Type
}

func (e *StorageInvalidValueError) Error() string {
	return fmt.Sprintf("%s is not supported value in this simple storage, use %s", e.Invalid, e.Valid)
}

// Validator is
type Validator struct {
}

func (s *Validator) ValidateKey(k interface{}) (string, error) {
	key, ok := k.(string)
	if !ok {
		return "", &StorageInvalidKeyError{
			Valid:   reflect.TypeOf((string)("")),
			Invalid: reflect.TypeOf(k),
		}
	}
	return key, nil
}

func (s *Validator) ValidateValue(v interface{}) ([]byte, error) {
	value, ok := v.([]byte)
	if !ok {
		return []byte{}, &StorageInvalidValueError{
			Valid:   reflect.TypeOf(([]byte)("")),
			Invalid: reflect.TypeOf(v),
		}
	}
	return value, nil
}
