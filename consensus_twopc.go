package transparent

import "github.com/juntaki/transparent/twopc"

// NewTwoPCConsensus returns Two phase commit consensus layer
func NewTwoPCConsensus() *Consensus {
	c := &Consensus{}
	participant := twopc.NewParticipant(c.commit)
	c.Participant = participant

	return c
}
