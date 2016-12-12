package twopc

import (
	"testing"

	"github.com/juntaki/transparent"
	"github.com/juntaki/transparent/dummy"
	test "github.com/juntaki/transparent/test"
)

func TestConsensus(t *testing.T) {
	//twopc.DebugLevel = 3
	var err error
	NewCoodinator()
	src1, err := dummy.NewSource(0)
	if err != nil {
		t.Error(err)
	}
	src2, err := dummy.NewSource(0)
	if err != nil {
		t.Error(err)
	}
	src3, err := dummy.NewSource(0)
	if err != nil {
		t.Error(err)
	}
	a1 := NewConsensus()
	a2 := NewConsensus()
	a3 := NewConsensus()

	s1 := transparent.NewStack()
	s2 := transparent.NewStack()
	s3 := transparent.NewStack()

	s1.Stack(src1)
	s2.Stack(src2)
	s3.Stack(src3)

	s1.Stack(a1)
	s2.Stack(a2)
	s3.Stack(a3)

	s1.Start()
	s2.Start()
	s3.Start()

	test.BasicStackFunc(t, s1)
	test.BasicStackFunc(t, s2)
	test.BasicStackFunc(t, s3)

	err = s1.Set("test1", "value1")
	if err != nil {
		t.Error(err)
	}
	err = s2.Set("test2", "value2")
	if err != nil {
		t.Error(err)
	}
	err = s3.Set("test3", "value3")
	if err != nil {
		t.Error(err)
	}

	err = s1.Sync()
	if err != nil {
		t.Error(err)
	}

	val, err := s1.Get("test1")
	if err != nil {
		t.Error(err)
	}
	if val != "value1" {
		t.Error(val)
	}

	err = s1.Remove("test1")
	if err != nil {
		t.Error(err)
	}

	val, err = s1.Get("test1")
	if err == nil {
		t.Error(err)
		t.Error(val)
	}

	s1.Stop()
	s2.Stop()
	s3.Stop()
}
