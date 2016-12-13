package custom

import (
	"testing"

	"github.com/juntaki/transparent/test"
)

func TestCustom(t *testing.T) {
	// Custom
	ds := test.NewStorage(0)
	cs, err := NewStorage(ds.Get, ds.Add, ds.Remove)
	if err != nil {
		t.Fatal(err)
	}
	test.BasicStorageFunc(t, cs)
}
