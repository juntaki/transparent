package transparent

import (
	"fmt"
	"reflect"
)

// Storage defines the interface that backend data storage destination should have.
// Add should not be failed.
type Storage interface {
	Get(key interface{}) (value interface{}, err error)
	Add(key interface{}, value interface{}) error
	Remove(key interface{}) error
}

// SimpleStorageInvalidKeyError means type of key is invalid for storage
type SimpleStorageInvalidKeyError struct {
	valid   reflect.Type
	invalid reflect.Type
}

func (e *SimpleStorageInvalidKeyError) Error() string {
	return fmt.Sprintf("%s is not supported key in the storage, use %s", e.invalid, e.valid)
}

// SimpleStorageInvalidValueError means type of value is invalid for storage
type SimpleStorageInvalidValueError struct {
	valid   reflect.Type
	invalid reflect.Type
}

func (e *SimpleStorageInvalidValueError) Error() string {
	return fmt.Sprintf("%s is not supported value in this simple storage, use %s", e.invalid, e.valid)
}

// StorageKeyNotFoundError means key is not found in the storage
type StorageKeyNotFoundError struct {
	Key interface{}
}

func (e *StorageKeyNotFoundError) Error() string { return "requested key is not found" }
