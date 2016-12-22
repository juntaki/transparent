package transparent

import (
	"errors"
	"sync"

	uuid "github.com/satori/go.uuid"
)

// NewLayerConsensus returns LayerConsensus.
// LayerConsensus wraps BackendTransmitter.
// It send Set operation and key-value to multiple Stacks asynchronously
// and Get key-value from Next Layer.
// It must be Stacked on a Layer.
//
//    User program A       User program B
//          |                    |
// -----------------------------------------
// |         transparent.Consensus         |
// -------------------- --------------------
// |transparent.Source| |transparent.Source|
// -------------------- --------------------
func NewLayerConsensus(t BackendTransmitter) (Layer, error) {
	c := &layerConsensus{
		inFlight:    make(map[string]chan error),
		Transmitter: t,
	}
	err := t.SetCallback(c.commit)
	if err != nil {
		return nil, err
	}
	return c, nil
}

type layerConsensus struct {
	lock        sync.Mutex
	inFlight    map[string]chan error
	next        Layer
	Transmitter BackendTransmitter
}

// Set send a request to cluster
func (d *layerConsensus) Set(key interface{}, value interface{}) (err error) {
	// We will check which message is commited by UUID
	uuid := uuid.NewV4().String()
	channel := make(chan error)
	d.lock.Lock()
	d.inFlight[uuid] = channel
	d.lock.Unlock()
	operation := &Message{
		Key:     key,
		Value:   value,
		Message: MessageSet,
		UUID:    uuid,
	}
	_, err = d.Transmitter.Request(operation)
	if err != nil {
		return err
	}
	err = <-channel
	d.lock.Lock()
	delete(d.inFlight, uuid)
	d.lock.Unlock()
	return err
}

// Get just get the value from next layer
func (d *layerConsensus) Get(key interface{}) (value interface{}, err error) {
	// Recursively get value from list.
	if d.next == nil {
		return nil, errors.New("next layer not found")
	}
	value, err = d.next.Get(key)
	if err != nil {
		return nil, err
	}
	return value, nil
}

// Remove send a request to cluster
func (d *layerConsensus) Remove(key interface{}) (err error) {
	uuid := uuid.NewV4().String()
	channel := make(chan error)
	d.lock.Lock()
	d.inFlight[uuid] = channel
	d.lock.Unlock()
	operation := &Message{
		Key:     key,
		Value:   nil,
		Message: MessageRemove,
		UUID:    uuid,
	}
	_, err = d.Transmitter.Request(operation)
	if err != nil {
		return err
	}
	err = <-channel
	d.lock.Lock()
	delete(d.inFlight, uuid)
	d.lock.Unlock()
	return err
}

// Sync send a request to cluster
func (d *layerConsensus) Sync() (err error) {
	uuid := uuid.NewV4().String()
	channel := make(chan error)
	d.lock.Lock()
	d.inFlight[uuid] = channel
	d.lock.Unlock()
	operation := &Message{
		Key:     nil,
		Value:   nil,
		Message: MessageSync,
		UUID:    uuid,
	}
	_, err = d.Transmitter.Request(operation)
	if err != nil {
		return err
	}
	err = <-channel
	d.lock.Lock()
	delete(d.inFlight, uuid)
	d.lock.Unlock()
	return err
}

// commit is callback function to apply operation
func (d *layerConsensus) commit(op *Message) (res *Message, err error) {
	err = nil
	key := op.Key
	if d.next == nil {
		err = errors.New("next layer not found")
	}
	switch op.Message {
	case MessageSync:
		err = d.next.Sync()
	case MessageRemove:
		err = d.next.Remove(key)
	case MessageSet:
		err = d.next.Set(key, op.Value)
	default:
		err = errors.New("unknown message")
	}
	d.lock.Lock()
	channel, ok := d.inFlight[op.UUID]
	d.lock.Unlock()
	if ok {
		channel <- err
	}
	return nil, err
}

func (d *layerConsensus) setNext(next Layer) error {
	d.next = next
	return nil
}

func (d *layerConsensus) start() error {
	return d.Transmitter.Start()
}

func (d *layerConsensus) stop() error {
	return d.Transmitter.Stop()
}
