package twopc

import "github.com/juntaki/transparent"

// NewConsensus returns Two phase commit consensus layer
func NewConsensus(serverAddr string) (*transparent.Consensus, error) {
	c := transparent.NewConsensus()
	participant, err := NewParticipant(serverAddr, c.Commit)
	if err != nil {
		return nil, err
	}
	c.Participant = participant

	return c, nil
}
