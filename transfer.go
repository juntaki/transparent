package transparent

import "errors"

// LayerReceiver wraps BackendReceiver.
// It receive operation and key-value from another Stack.
// This layer must be the top of Stack.
type LayerReceiver struct {
	Receiver BackendReceiver
}

// NewLayerReceiver returns LayerReceiver implements Layer
func NewLayerReceiver(Receiver BackendReceiver) *LayerReceiver {
	return &LayerReceiver{
		Receiver: Receiver,
	}
}

// Set is not allowed, operation should be transfered from Transmitter.
func (r *LayerReceiver) Set(key interface{}, value interface{}) error {
	return errors.New("don't send Set")
}

// Get is not allowed, operation should be transfered from Transmitter.
func (r *LayerReceiver) Get(key interface{}) (value interface{}, err error) {
	return nil, errors.New("don't send Get")
}

// Remove is not allowed, operation should be transfered from Transmitter.
func (r *LayerReceiver) Remove(key interface{}) error {
	return errors.New("don't send Remove")
}

// Sync is not allowed, operation should be transfered from Transmitter.
func (r *LayerReceiver) Sync() error {
	return errors.New("don't send Sync")
}

func (r *LayerReceiver) setNext(l Layer) error {
	return r.Receiver.SetNext(l)
}
func (r *LayerReceiver) start() error {
	return r.Receiver.Start()
}
func (r *LayerReceiver) stop() error {
	return r.Receiver.Stop()
}

// LayerTransmitter wraps BackendTransmitter.
// It send operation and key-value to another Stack.
// This layer must be the bottom of Stack.
type LayerTransmitter struct {
	Transmitter BackendTransmitter
}

// NewLayerTransmitter returns LayerTransmitter implements Layer
func NewLayerTransmitter(Transmitter BackendTransmitter) *LayerTransmitter {
	return &LayerTransmitter{
		Transmitter: Transmitter,
	}
}

// Set convert key-value to Message and Request it.
func (r *LayerTransmitter) Set(key interface{}, value interface{}) error {
	operation := &Message{
		Message: MessageSet,
		Key:     key,
		Value:   value,
	}
	_, err := r.Transmitter.Request(operation)
	if err != nil {
		return err
	}
	return nil
}

// Get convert key to Message and Request it.
func (r *LayerTransmitter) Get(key interface{}) (value interface{}, err error) {
	operation := &Message{
		Message: MessageGet,
		Key:     key,
	}

	feature, err := r.Transmitter.Request(operation)
	if err != nil {
		return nil, err
	}

	return feature.Value, nil
}

// Remove convert key to Message and Request it.
func (r *LayerTransmitter) Remove(key interface{}) error {
	operation := &Message{
		Message: MessageRemove,
		Key:     key,
	}
	_, err := r.Transmitter.Request(operation)
	if err != nil {
		return err
	}
	return nil
}

// Sync makes Message and Request it.
func (r *LayerTransmitter) Sync() error {
	operation := &Message{
		Message: MessageSync,
	}
	_, err := r.Transmitter.Request(operation)
	if err != nil {
		return err
	}
	return nil
}
func (r *LayerTransmitter) setNext(l Layer) error {
	return errors.New("don't send next layer")
}
func (r *LayerTransmitter) start() error {
	return r.Transmitter.Start()
}
func (r *LayerTransmitter) stop() error {
	return r.Transmitter.Stop()
}
