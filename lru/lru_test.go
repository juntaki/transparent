package lru

import "testing"

func TestLRUcache(t *testing.T) {
	c := New(10)
	val, err := c.Get("test")
	if err == nil {
		t.Error(val, err)
	}
	err = c.Add("test", "value")
	if err != nil {
		t.Error(err)
	}
	val, err = c.Get("test")
	if err != nil {
		t.Error(val, err)
	}
	if val != "value" {
		t.Error(val, err)
	}
	err = c.Remove("test")
	if err != nil {
		t.Error(err)
	}
	val, err = c.Get("test")
	if err == nil {
		t.Error(val, err)
	}
}
