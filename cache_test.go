package transparent

import "testing"

func basicCacheFunc(t *testing.T, c *Cache) {
	s, err := NewDummySource(0)
	if err != nil {
		t.Error(err)
	}
	Stack(c, s)
	basicLayerFunc(t, c)
}

func TestCache(t *testing.T) {
	var c *Cache
	var err error
	c, err = NewLRUCache(10, 100)
	if err != nil {
		t.Error(err)
	}
	basicCacheFunc(t, c)

	c, err = NewFilesystemCache(10, "/tmp")
	if err != nil {
		t.Error(err)
	}
	basicCacheFunc(t, c)

	svc, err := newMockS3Client()
	if err != nil {
		t.Error(err)
	}
	c, err = NewS3Cache(10, "bucket", svc)
	if err != nil {
		t.Error(err)
	}
	basicCacheFunc(t, c)
}
