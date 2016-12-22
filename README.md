# transparent

Transparent cache and distributed commit to key-value store written in Go.

[![Build Status](https://travis-ci.org/juntaki/transparent.svg?branch=master)](https://travis-ci.org/juntaki/transparent)
[![Coverage Status](https://coveralls.io/repos/github/juntaki/transparent/badge.svg?branch=master)](https://coveralls.io/github/juntaki/transparent?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/juntaki/transparent)](https://goreportcard.com/report/github.com/juntaki/transparent)
[![GoDoc](https://godoc.org/github.com/juntaki/transparent?status.svg)](https://godoc.org/github.com/juntaki/transparent)


## Documentation

### Installation

~~~ sh
go get github.com/juntaki/transparent
~~~

### Basic usage

First, create layers and stack them with Stack.Stack().
This example adds LRU memory cache and filesystem cache to dummy source layer.

~~~go
	cacheLayer1, _ := lru.NewCache(10, 100)
	cacheLayer2 := filesystem.NewCache(10, "/tmp")
	sourceLayer := test.NewSource(10)

	stack := transparent.NewStack()
	stack.Stack(sourceLayer)
	stack.Stack(cacheLayer2)
	stack.Stack(cacheLayer1)
~~~

If you manipulate the Stack, the value will be transmitted from top layer to the bottom layer transparently.

~~~go
	stack.Set("key", []byte("value"))
	stack.Sync()
    
    // value, _ = cacheLayer1.Get("key") // "value"
	// value, _ = cacheLayer2.Get("key") // "value"
	// value, _ = sourceLayer.Get("key") // "value"

	value, _ := stack.Get("key")
	fmt.Printf("%s\n", value)            // "value"
~~~

For details, please refer to [Godoc] (https://godoc.org/github.com/juntaki/transparent).
