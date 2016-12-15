package twopc

import "github.com/juntaki/transparent"

// NewConsensus returns Two phase commit consensus layer
func NewConsensus(serverAddr string) (*transparent.LayerConsensus, error) {
	c := transparent.NewLayerConsensus()
	participant, err := NewParticipant(serverAddr, c.Commit)
	if err != nil {
		return nil, err
	}
	c.Transmitter = participant

	return c, nil
}
