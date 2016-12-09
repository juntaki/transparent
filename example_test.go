package transparent_test

import (
	"fmt"

	"github.com/juntaki/transparent"
)

func Example() {
	var err error
	cacheLayer1, _ := transparent.NewLRUCache(10, 100)
	cacheLayer2, _ := transparent.NewFilesystemCache(10, "/tmp")
	sourceLayer, _ := transparent.NewDummySource(10)
	transparent.Stack(cacheLayer1, cacheLayer2)
	transparent.Stack(cacheLayer2, sourceLayer)

	cacheLayer1.Set("key", []byte("value"))
	cacheLayer1.Sync()
	value, _ := cacheLayer1.Get("key")
	fmt.Printf("%s\n", value)

	value, err = cacheLayer2.Get("key")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", value)
	value, err = sourceLayer.Get("key")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", value)
	// Output:
	// value
	// value
	// value
}
