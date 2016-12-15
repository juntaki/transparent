package transparent

import "errors"

type Receiver interface {
	Start() error
	Stop() error
	SetNext(l Layer) error
}

type LayerReceiver struct {
	Receiver Receiver
}

// NewLayerReceiver returns LayerReceiver
func NewLayerReceiver(Receiver Receiver) *LayerReceiver {
	return &LayerReceiver{
		Receiver: Receiver,
	}
}

func (r *LayerReceiver) Set(key interface{}, value interface{}) error {
	return errors.New("don't send Set")
}
func (r *LayerReceiver) Get(key interface{}) (value interface{}, err error) {
	return nil, errors.New("don't send Get")
}
func (r *LayerReceiver) Remove(key interface{}) error {
	return errors.New("don't send Remove")
}
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

// Transmitter is interface to another system
type Transmitter interface {
	Request(operation *Message) (*Message, error)
	Start() error
	Stop() error
}

type LayerTransmitter struct {
	Transmitter Transmitter
}

// NewLayerTransmitter returns LayerTransmitter
func NewLayerTransmitter(Transmitter Transmitter) *LayerTransmitter {
	return &LayerTransmitter{
		Transmitter: Transmitter,
	}
}

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
