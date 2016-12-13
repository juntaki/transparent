package transparent

// Storage defines the interface that backend data storage destination should have.
// Add should not be failed.
type Storage interface {
	Get(key interface{}) (value interface{}, err error)
	Add(key interface{}, value interface{}) error
	Remove(key interface{}) error
}
