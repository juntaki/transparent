package transparent

import "errors"

type layerReceiver struct {
	Receiver BackendReceiver
}

// NewLayerReceiver returns LayerReceiver.
// LayerReceiver wraps BackendReceiver.
// It receive operation and key-value from another Stack.
// This layer must be the top of Stack.
func NewLayerReceiver(Receiver BackendReceiver) Layer {
	return &layerReceiver{
		Receiver: Receiver,
	}
}

// Set is not allowed, operation should be transfered from Transmitter.
func (r *layerReceiver) Set(key interface{}, value interface{}) error {
	return errors.New("don't send Set")
}

// Get is not allowed, operation should be transfered from Transmitter.
func (r *layerReceiver) Get(key interface{}) (value interface{}, err error) {
	return nil, errors.New("don't send Get")
}

// Remove is not allowed, operation should be transfered from Transmitter.
func (r *layerReceiver) Remove(key interface{}) error {
	return errors.New("don't send Remove")
}

// Sync is not allowed, operation should be transfered from Transmitter.
func (r *layerReceiver) Sync() error {
	return errors.New("don't send Sync")
}

func (r *layerReceiver) setNext(l Layer) error {
	return r.Receiver.SetNext(l)
}
func (r *layerReceiver) start() error {
	return r.Receiver.Start()
}
func (r *layerReceiver) stop() error {
	return r.Receiver.Stop()
}

type layerTransmitter struct {
	Transmitter BackendTransmitter
}

// NewLayerTransmitter returns LayerTransmitter.
// LayerTransmitter wraps BackendTransmitter.
// It send operation and key-value to another Stack.
// This layer must be the bottom of Stack.
func NewLayerTransmitter(Transmitter BackendTransmitter) Layer {
	return &layerTransmitter{
		Transmitter: Transmitter,
	}
}

// Set convert key-value to Message and Request it.
func (r *layerTransmitter) Set(key interface{}, value interface{}) error {
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
func (r *layerTransmitter) Get(key interface{}) (value interface{}, err error) {
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
func (r *layerTransmitter) Remove(key interface{}) error {
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
func (r *layerTransmitter) Sync() error {
	operation := &Message{
		Message: MessageSync,
	}
	_, err := r.Transmitter.Request(operation)
	if err != nil {
		return err
	}
	return nil
}
func (r *layerTransmitter) setNext(l Layer) error {
	return errors.New("don't send next layer")
}
func (r *layerTransmitter) start() error {
	return r.Transmitter.Start()
}
func (r *layerTransmitter) stop() error {
	return r.Transmitter.Stop()
}
