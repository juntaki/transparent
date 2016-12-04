package twopc

import (
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	DebugLevel = 3
	NewCoodinator()

	array1 := [][]interface{}{}
	commit1 := func(key, value interface{}) error {
		array1 = append(array1, []interface{}{key, value})
		return nil
	}
	array2 := [][]interface{}{}
	commit2 := func(key, value interface{}) error {
		array2 = append(array2, []interface{}{key, value})
		return nil
	}
	array3 := [][]interface{}{}
	commit3 := func(key, value interface{}) error {
		array3 = append(array3, []interface{}{key, value})
		return nil
	}

	expected := [][]interface{}{
		{"testkey1", "testvalue1"},
		{"testkey2", "testvalue2"},
		{"testkey3", "testvalue3"},
	}
	a1 := NewParticipant(commit1)
	a2 := NewParticipant(commit2)
	a3 := NewParticipant(commit3)

	a1.Request("testkey1", "testvalue1")
	a2.Request("testkey2", "testvalue2")
	a3.Request("testkey3", "testvalue3")

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
}

func TestEncode(t *testing.T) {
	kv := &keyValue{
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
