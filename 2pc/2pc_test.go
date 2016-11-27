package twopc

import (
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	DebugLevel = 5
	c := Coodinator{}
	go c.StartServ(1000)

	a := Attendee{}
	go a.StartClient(1000)
	a2 := Attendee{}
	go a2.StartClient(1000)
	a3 := Attendee{}
	go a3.StartClient(1000)

	time.Sleep(5 * time.Second)
	a.Set()
	a.Set()
	time.Sleep(5 * time.Second)
}
