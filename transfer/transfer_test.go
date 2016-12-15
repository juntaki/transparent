package transfer

import (
	"testing"

	"github.com/juntaki/transparent"
	"github.com/juntaki/transparent/test"
)

func TestTransfer(t *testing.T) {
	serverAddr := "localhost:8080"
	r := NewReceiver(serverAddr)
	d := test.NewSource(0)
	r.SetNext(d)
	r.Start()

	a1 := NewSimpleTransmitter(serverAddr)
	tra := transparent.NewLayerTransmitter(a1)

	test.BasicTransmitterFunc(t, tra)
}
