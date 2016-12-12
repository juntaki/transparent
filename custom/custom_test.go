package custom

import (
	"testing"

	"github.com/juntaki/transparent/dummy"
	test "github.com/juntaki/transparent/test"
)

func TestCustom(t *testing.T) {
	var err error

	// Custom
	ds, err := dummy.NewStorage(0)
	if err != nil {
		t.Fatal(err)
	}
	cs, err := NewStorage(ds.Get, ds.Add, ds.Remove)
	if err != nil {
		t.Fatal(err)
	}
	test.BasicStorageFunc(t, cs)
}
