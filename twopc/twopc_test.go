package twopc

import (
	"testing"
	"time"

	"github.com/juntaki/transparent"
)

func TestServer(t *testing.T) {
	DebugLevel = 3
	serverAddr := "localhost:8888"
	NewCoodinator(serverAddr)

	array1 := [][]interface{}{}
	commit1 := func(op *transparent.Message) (*transparent.Message, error) {
		array1 = append(array1, []interface{}{op.Key, op.Value})
		return nil, nil
	}
	array2 := [][]interface{}{}
	commit2 := func(op *transparent.Message) (*transparent.Message, error) {
		array2 = append(array2, []interface{}{op.Key, op.Value})
		return nil, nil
	}
	array3 := [][]interface{}{}
	commit3 := func(op *transparent.Message) (*transparent.Message, error) {
		array3 = append(array3, []interface{}{op.Key, op.Value})
		return nil, nil
	}

	expected := [][]interface{}{
		{"testkey1", "testvalue1"},
		{"testkey2", "testvalue2"},
		{"testkey3", "testvalue3"},
	}
	a1 := NewParticipant(serverAddr)
	a2 := NewParticipant(serverAddr)
	a3 := NewParticipant(serverAddr)

	err := a1.SetCallback(commit1)
	if err != nil {
		t.Fatal(err)
	}
	err = a2.SetCallback(commit2)
	if err != nil {
		t.Fatal(err)
	}
	err = a3.SetCallback(commit3)
	if err != nil {
		t.Fatal(err)
	}

	a1.Start()
	a2.Start()
	a3.Start()

	a1.Request(&transparent.Message{Key: "testkey1", Value: "testvalue1"})
	a2.Request(&transparent.Message{Key: "testkey2", Value: "testvalue2"})
	a3.Request(&transparent.Message{Key: "testkey3", Value: "testvalue3"})

	time.Sleep(3 * time.Second)

	if len(expected) != len(array1) ||
		len(expected) != len(array2) ||
		len(expected) != len(array3) {
		t.Error(array1, array2, array3)
		return
	}
	for i := range expected {
		if expected[i][0] != array1[i][0] ||
			expected[i][1] != array1[i][1] ||
			expected[i][0] != array2[i][0] ||
			expected[i][1] != array2[i][1] ||
			expected[i][0] != array3[i][0] ||
			expected[i][1] != array3[i][1] {
			t.Error(array1, array2, array3)
		}
	}
	a1.Stop()
	a2.Stop()
	a3.Stop()
}

func TestEncode(t *testing.T) {
	kv := &transparent.Message{
		Key:   "testKey",
		Value: "testValue",
	}

	a := Participant{}
	req, err := a.encode(kv)
	if err != nil {
		t.Error(err)
	}

	kv2, err := a.decode(req.Payload)
	if err != nil {
		t.Error(err)
	}

	if kv == kv2 {
		t.Error(kv, kv2)
	}

}
