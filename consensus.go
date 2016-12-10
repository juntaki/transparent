package transparent

import (
	"errors"
	"sync"

	"github.com/juntaki/transparent/twopc"
	uuid "github.com/satori/go.uuid"
)

//    User program A       User program B
//          |                    |
// -----------------------------------------
// |         transparent.Consensus         |
// -------------------- --------------------
// |transparent.Source| |transparent.Source|
// -------------------- --------------------

// Participant is interface to consensus algorithm
type Participant interface {
	Request(key interface{}, value interface{}) error
}

// Consensus layer provide transactional write to cluster.
// There is no storage, the layer of only forward.
type Consensus struct {
	lock        sync.Mutex
	inFlight    map[string]chan error
	upper       Layer
	lower       Layer
	Participant Participant
}

// NewTwoPCConsensus returns Two phase commit consensus layer
func NewTwoPCConsensus() *Consensus {
	c := &Consensus{inFlight: make(map[string]chan error)}
	participant := twopc.NewParticipant(c.commit)
	c.Participant = participant

	return c
}

// Set send a request to cluster
func (d *Consensus) Set(key interface{}, value interface{}) (err error) {
	uuid := uuid.NewV4().String()
	channel := make(chan error)
	d.lock.Lock()
	d.inFlight[uuid] = channel
	d.lock.Unlock()
	operation := operation{
		Value:   value,
		Message: messageSet,
		UUID:    uuid,
	}
	err = d.Participant.Request(key, operation)
	if err != nil {
		return err
	}
	err = <-channel
	d.lock.Lock()
	delete(d.inFlight, uuid)
	d.lock.Unlock()
	return err
}

// Get just get the value from lower layer
func (d *Consensus) Get(key interface{}) (value interface{}, err error) {
	// Recursively get value from list.
	if d.lower == nil {
		return nil, errors.New("lower layer not found")
	}
	value, err = d.lower.Get(key)
	if err != nil {
		return nil, err
	}
	return value, nil
}

// Remove send a request to cluster
func (d *Consensus) Remove(key interface{}) (err error) {
	uuid := uuid.NewV4().String()
	channel := make(chan error)
	d.lock.Lock()
	d.inFlight[uuid] = channel
	d.lock.Unlock()
	operation := operation{
		Value:   nil,
		Message: messageRemove,
		UUID:    uuid,
	}
	err = d.Participant.Request(key, operation)
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
func (d *Consensus) Sync() (err error) {
	uuid := uuid.NewV4().String()
	channel := make(chan error)
	d.lock.Lock()
	d.inFlight[uuid] = channel
	d.lock.Unlock()
	operation := operation{
		Value:   nil,
		Message: messageSync,
		UUID:    uuid,
	}
	err = d.Participant.Request(nil, operation)
	if err != nil {
		return err
	}
	err = <-channel
	d.lock.Lock()
	delete(d.inFlight, uuid)
	d.lock.Unlock()
	return err
}

// commit should be callback function of message receiver
func (d *Consensus) commit(key interface{}, value interface{}) (err error) {
	err = nil
	operation, ok := value.(operation)
	if !ok {
		return errors.New("value should be operation")
	}
	if d.lower == nil {
		err = errors.New("lower layer not found")
	}
	switch operation.Message {
	case messageSync:
		err = d.lower.Sync()
	case messageRemove:
		err = d.lower.Remove(key)
	case messageSet:
		err = d.lower.Set(key, operation.Value)
	default:
		err = errors.New("unknown message")
	}
	d.lock.Lock()
	channel, ok := d.inFlight[operation.UUID]
	d.lock.Unlock()
	if ok {
		channel <- err
	}
	return err
}

// SetUpper set upper layer
func (d *Consensus) setUpper(upper Layer) {
	d.upper = upper
}

// SetLower set lower layer
func (d *Consensus) setLower(lower Layer) {
	d.lower = lower
}
