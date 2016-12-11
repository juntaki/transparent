package transparent

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"

	"github.com/pkg/errors"
)

type simpleStorageWrapper struct {
	Storage
}

// Get is file read
func (f *simpleStorageWrapper) Get(k interface{}) (interface{}, error) {
	key, err := f.encodeKey(k)
	if err != nil {
		return nil, err
	}
	v, err := f.Storage.Get(key)
	if err != nil {
		_, ok := err.(*StorageKeyNotFoundError)
		if ok {
			return nil, &StorageKeyNotFoundError{Key: k}
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
func (f *simpleStorageWrapper) Add(k interface{}, v interface{}) error {
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
func (f *simpleStorageWrapper) Remove(k interface{}) error {
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

func (f *simpleStorageWrapper) encodeKey(k interface{}) (string, error) {
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

func (f *simpleStorageWrapper) encodeValue(v interface{}) ([]byte, error) {
	gob.Register(v)
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	cause := encoder.Encode(&v)
	if cause != nil {
		return nil, errors.Wrap(cause, "failed to encode value")
	}
	return buf.Bytes(), nil
}

func (f *simpleStorageWrapper) decodeValue(v []byte) (interface{}, error) {
	var d interface{}
	buf := bytes.NewBuffer(v)
	decoder := gob.NewDecoder(buf)
	cause := decoder.Decode(&d)
	if cause != nil {
		return nil, errors.Wrap(cause, "failed to decode value")
	}
	return d, nil
}
