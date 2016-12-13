package filesystem

import "github.com/juntaki/transparent"

// NewCache returns FilesystemCache
func NewCache(bufferSize int, directory string) *transparent.Cache {
	filesystem := NewSimpleStorage(directory)
	layer, _ := transparent.NewCache(bufferSize, filesystem)
	return layer
}

// NewSource returns FilesystemSource
func NewSource(directory string) *transparent.Source {
	filesystem := NewSimpleStorage(directory)
	layer, _ := transparent.NewSource(filesystem)
	return layer
}
