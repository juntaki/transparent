package twopc

import (
	"testing"

	"github.com/juntaki/transparent/test"
)

func TestConsensus(t *testing.T) {
	DebugLevel = 3
	serverAddr := "localhost:8080"
	_, err := NewCoodinator(serverAddr)
	if err != nil {
		t.Fatal(err)
	}

	a1, err := NewConsensus(serverAddr)
	if err != nil {
		t.Fatal(err)
	}
	a2, err := NewConsensus(serverAddr)
	if err != nil {
		t.Fatal(err)
	}

	test.BasicConsensusFunc(t, a1, a2)
}
