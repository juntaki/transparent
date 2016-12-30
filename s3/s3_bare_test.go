package s3

import (
	"bytes"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
)

func TestS3Bare(t *testing.T) {
	b := NewBare()
	b.Value["Key"] = aws.String("key")
	b.Value["Bucket"] = aws.String("bucket")
	b.Value["Body"] = bytes.NewReader([]byte("test"))
	err := b.set()
	if err != nil {
		t.Error(err)
	}
	b.get(b.getObjectInput)
	b.get(b.getObjectOutput)
	b.get(b.putObjectInput)
	b.get(b.deleteObjectInput)
}
