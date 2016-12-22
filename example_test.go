package transparent_test

import (
	"fmt"

	"github.com/juntaki/transparent"
	"github.com/juntaki/transparent/filesystem"
	"github.com/juntaki/transparent/lru"
	"github.com/juntaki/transparent/test"
)

func Example() {
	var err error
	cacheLayer1, _ := lru.NewCache(10, 100)
	cacheLayer2 := filesystem.NewCache(10, "/tmp")
	sourceLayer := test.NewSource(10)

	stack := transparent.NewStack()
	stack.Stack(sourceLayer)
	stack.Stack(cacheLayer2)
	stack.Stack(cacheLayer1)

	stack.Start()
	defer stack.Stop()

	stack.Set("key", []byte("value"))
	stack.Sync()

	value, _ := stack.Get("key")
	fmt.Printf("%s\n", value)

	value, _ = cacheLayer1.Get("key")
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
	// value
}
