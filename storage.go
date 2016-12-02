package transparent

// Storage defines the interface that backend data storage destination should have.
// Add should not be failed.
type Storage interface {
	Get(key interface{}) (value interface{}, found bool)
	Add(key interface{}, value interface{})
	Remove(key interface{})
}
