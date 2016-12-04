package transparent

import (
	"testing"
	"time"

	"github.com/juntaki/transparent/twopc"
)

func TestConsensus(t *testing.T) {
	twopc.DebugLevel = 3
	var err error
	twopc.NewCoodinator()
	src1 := NewDummySource(5)
	src2 := NewDummySource(5)
	src3 := NewDummySource(5)
	a1 := NewTwoPCConsensus()
	a2 := NewTwoPCConsensus()
	a3 := NewTwoPCConsensus()
	Stack(a1, src1)
	Stack(a2, src2)
	Stack(a3, src3)

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

	time.Sleep(3 * time.Second)

	val, err := a1.Get("test1")
	if err != nil {
		t.Error(err)
	}
	if val != "value1" {
		t.Error(val)
	}
}
