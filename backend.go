package transparent

// BackendReceiver is interface from another system
// Operation should apply to Next layer
type BackendReceiver interface {
	Start() error
	Stop() error
	SetNext(l Layer) error
}

// BackendTransmitter is interface to another system
// Request transfer an operation as Message
// If request will be processed asynchronously,
// callback function should executed with reply Message.
type BackendTransmitter interface {
	Request(operation *Message) (*Message, error)
	Start() error
	Stop() error
	SetCallback(func(m *Message) error) error
}

// BackendStorage defines the interface that backend data storage.
type BackendStorage interface {
	Get(key interface{}) (value interface{}, err error)
	Add(key interface{}, value interface{}) error
	Remove(key interface{}) error
}
