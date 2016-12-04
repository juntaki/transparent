package transparent

import "errors"

//    User program A       User program B
//          |                    |
// -----------------------------------------
// |         transparent.Consensus         |
// -------------------- --------------------
// |transparent.Source| |transparent.Source|
// -------------------- --------------------

// Participant is interface to consensus algorhytm
type Participant interface {
	Request(key interface{}, value interface{}) error
}

// Consensus layer provide transactional write to cluster.
// There is no storage, the layer of only forward.
type Consensus struct {
	upper       Layer
	lower       Layer
	Participant Participant
}

// Set send a request to cluster
func (d *Consensus) Set(key interface{}, value interface{}) (err error) {
	operation := operation{
		Value:   value,
		Message: messageSet,
	}
	return d.Participant.Request(key, operation)
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
	operation := operation{
		Value:   nil,
		Message: messageRemove,
	}
	return d.Participant.Request(key, operation)
}

// Sync send a request to cluster
func (d *Consensus) Sync() (err error) {
	operation := operation{
		Value:   nil,
		Message: messageSync,
	}
	return d.Participant.Request(nil, operation)
}

// commit should be callback function of message receiver
func (d *Consensus) commit(key interface{}, value interface{}) error {
	operation, ok := value.(operation)
	if !ok {
		return errors.New("value should be operation")
	}
	if d.lower == nil {
		return errors.New("lower layer not found")
	}
	switch operation.Message {
	case messageSync:
		return d.lower.Sync()
	case messageRemove:
		return d.lower.Remove(key)
	case messageSet:
		return d.lower.Set(key, operation.Value)
	}
	return errors.New("unknown message")
}

// SetUpper set upper layer
func (d *Consensus) setUpper(upper Layer) {
	d.upper = upper
}

// SetLower set lower layer
func (d *Consensus) setLower(lower Layer) {
	d.lower = lower
}
