package filesystem

import "github.com/juntaki/transparent"

// NewCache returns FilesystemCache
func NewCache(bufferSize int, directory string) transparent.Layer {
	filesystem := NewSimpleStorage(directory)
	layer, _ := transparent.NewLayerCache(bufferSize, filesystem)
	return layer
}

// NewSource returns FilesystemSource
func NewSource(directory string) transparent.Layer {
	filesystem := NewSimpleStorage(directory)
	layer, _ := transparent.NewLayerSource(filesystem)
	return layer
}
