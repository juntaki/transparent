package transparent

import (
	"testing"

	"github.com/juntaki/transparent/twopc"
)

func TestConsensus(t *testing.T) {
	//twopc.DebugLevel = 3
	var err error
	twopc.NewCoodinator()
	src1, err := NewDummySource(0)
	if err != nil {
		t.Error(err)
	}
	src2, err := NewDummySource(0)
	if err != nil {
		t.Error(err)
	}
	src3, err := NewDummySource(0)
	if err != nil {
		t.Error(err)
	}
	a1 := NewTwoPCConsensus()
	a2 := NewTwoPCConsensus()
	a3 := NewTwoPCConsensus()
	Stack(a1, src1)
	Stack(a2, src2)
	Stack(a3, src3)

	basicLayerFunc(t, a1)
	basicLayerFunc(t, a2)
	basicLayerFunc(t, a3)

	err = a1.Set("test1", "value1")
	if err != nil {
		t.Error(err)
	}
	err = a2.Set("test2", "value2")
	if err != nil {
		t.Error(err)
	}
	err = a3.Set("test3", "value3")
	if err != nil {
		t.Error(err)
	}

	err = a1.Sync()
	if err != nil {
		t.Error(err)
	}

	val, err := a1.Get("test1")
	if err != nil {
		t.Error(err)
	}
	if val != "value1" {
		t.Error(val)
	}

	err = a1.Remove("test1")
	if err != nil {
		t.Error(err)
	}

	val, err = a1.Get("test1")
	if err == nil {
		t.Error(err)
		t.Error(val)
	}
}
