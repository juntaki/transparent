package transfer

import (
	"testing"

	"github.com/juntaki/transparent"
	"github.com/juntaki/transparent/test"
)

func TestTransfer(t *testing.T) {
	serverAddr := "localhost:8080"
	r := NewSimpleLayerReceiver(serverAddr)
	d := test.NewSource(0)
	s := transparent.NewStack()
	s.Stack(d)
	s.Stack(r)
	s.Start()

	tra := NewSimpleLayerTransmitter(serverAddr)

	test.BasicTransmitterFunc(t, tra)
	s.Stop()
}
