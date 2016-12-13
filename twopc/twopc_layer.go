package twopc

import "github.com/juntaki/transparent"

// NewConsensus returns Two phase commit consensus layer
func NewConsensus(serverAddr string) *transparent.Consensus {
	c := transparent.NewConsensus()
	participant := NewParticipant(serverAddr, c.Commit)
	c.Participant = participant

	return c
}
