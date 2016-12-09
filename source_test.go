package transparent

import "testing"

func TestSource(t *testing.T) {
	var err error
	var l Layer
	l, err = NewDummySource(1)
	if err != nil {
		t.Error(err)
	}
	basicLayerFunc(t, l)

	l, err = NewFilesystemSource("/tmp")
	if err != nil {
		t.Error(err)
	}
	basicLayerFunc(t, l)

	svc, err := newMockS3Client()
	if err != nil {
		t.Error(err)
	}
	l, err = NewS3Source("bucket", svc)
	if err != nil {
		t.Error(err)
	}
	basicLayerFunc(t, l)
}
