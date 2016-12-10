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

First, create layers and stack them with Stack().

~~~go
var err error
cacheLayer1, _ := transparent.NewLRUCache(10, 100)
cacheLayer2, _ := transparent.NewFilesystemCache(10, "/tmp")
sourceLayer := transparent.NewDummySource(10)
transparent.Stack(cacheLayer1, cacheLayer2)
transparent.Stack(cacheLayer2, sourceLayer)
~~~

If you manipulate the layer stacked on top, the value will be transmitted to the bottom layer.

~~~go
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
~~~

~~~go:result
value
value
value
~~~

For details, please refer to [Godoc] (https://godoc.org/github.com/juntaki/transparent).
