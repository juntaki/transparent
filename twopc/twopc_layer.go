package twopc

import "github.com/juntaki/transparent"

// NewTwoPCConsensus returns Two phase commit consensus layer
func NewConsensus() *transparent.Consensus {
	c := transparent.NewConsensus()
	participant := NewParticipant(c.Commit)
	c.Participant = participant

	return c
}
