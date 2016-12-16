package twopc

import "github.com/juntaki/transparent"

// NewConsensus returns Two phase commit consensus layer
func NewConsensus(serverAddr string) (*transparent.LayerConsensus, error) {
	participant := NewParticipant(serverAddr)
	c, err := transparent.NewLayerConsensus(participant)
	if err != nil {
		return nil, err
	}

	return c, nil
}
