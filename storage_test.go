package transparent

import "testing"

func TestCustomStorage(t *testing.T) {
	test := make(map[interface{}]interface{})

	getFunc := func(k interface{}) (interface{}, bool) {
		value, ok := test[k]
		return value, ok
	}
	addFunc := func(k interface{}, v interface{}) {
		test[k] = v
	}
	removeFunc := func(k interface{}) {
		delete(test, k)
	}

	storage, err := NewCustomStorage(getFunc, addFunc, removeFunc)
	if err != nil {
		t.Error(err)
	}
	storage.Add("test", "value")
	value, ok := storage.Get("test")
	if ok != true || value != "value" {
		t.Error(ok)
		t.Error(value)
	}
	storage.Remove("test")
	value2, ok := storage.Get("test")
	if ok != false {
		t.Error(ok)
		t.Error(value2)
	}
}
