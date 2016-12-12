package filesystem

import "github.com/juntaki/transparent"

// NewFilesystemCache returns FilesystemCache
func NewCache(bufferSize int, directory string) (*transparent.Cache, error) {
	filesystem, err := NewSimpleStorage(directory)
	if err != nil {
		return nil, err
	}
	layer, err := transparent.NewCacheLayer(bufferSize, filesystem)
	if err != nil {
		return nil, err
	}
	return layer, nil
}

// NewFilesystemSource returns FilesystemSource
func NewSource(directory string) (*transparent.Source, error) {
	filesystem, err := NewSimpleStorage(directory)
	if err != nil {
		return nil, err
	}
	layer, err := transparent.NewSource(filesystem)
	if err != nil {
		return nil, err
	}
	return layer, nil
}
