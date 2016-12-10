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

// StorageInvalidKeyError means type of key is invalid for storage
type StorageInvalidKeyError struct {
	valid   reflect.Type
	invalid reflect.Type
}

func (e *StorageInvalidKeyError) Error() string {
	return fmt.Sprintf("%s is not supported key in the storage, use %s", e.invalid, e.valid)
}

// StorageInvalidValueError means type of value is invalid for storage
type StorageInvalidValueError struct {
	valid   reflect.Type
	invalid reflect.Type
}

func (e *StorageInvalidValueError) Error() string {
	return fmt.Sprintf("%s is not supported value in this simple storage, use %s", e.invalid, e.valid)
}

// StorageKeyNotFoundError means key is not found in the storage
type StorageKeyNotFoundError struct {
	Key interface{}
}

func (e *StorageKeyNotFoundError) Error() string { return "requested key is not found" }
