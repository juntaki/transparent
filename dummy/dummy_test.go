package dummy

import "testing"

func TestDummy(t *testing.T) {
	var err error

	//
	_, err = NewStorage(0)
	if err != nil {
		t.Fatal(err)
	}
	//transparent_test.basicStorageFunc(t, ds)

	_, err = NewSource(1)
	if err != nil {
		t.Error(err)
	}
	//basicLayerFunc(t, l)
}
