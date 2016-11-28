package twopc

import (
	"fmt"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	DebugLevel = 3
	c := Coodinator{}
	go c.StartServ(1000)

	commit := func(key, value interface{}) {
		fmt.Println("key", key, "value", value)
	}
	a := Attendee{}
	go a.StartClient(1000, commit)
	a2 := Attendee{}
	go a2.StartClient(1000, commit)
	a3 := Attendee{}
	go a3.StartClient(1000, commit)

	time.Sleep(5 * time.Second)
	a.Set("testkey", "testvalue")
	a.Set("testkey", "testvalue")
	time.Sleep(5 * time.Second)
}

func TestEncode(t *testing.T) {
	kv := &keyValue{
		Key:   "testKey",
		Value: "testValue",
	}

	a := Attendee{}
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
