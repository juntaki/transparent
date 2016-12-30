package s3

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
)

func TestS3Bare(t *testing.T) {
	key := &BareKey{
		Key:    aws.String("test"),
		Bucket: aws.String("bucket"),
	}
	value := NewBareValue()
	sb := newBare(key, value)
	sb.Value["Key"] = aws.String("test2")
	err := sb.Set()
	if err != nil {
		t.Error(err)
	}
	sb.Get(sb.getObjectInput)
	sb.Get(sb.getObjectOutput)
	sb.Get(sb.putObjectInput)
	sb.Get(sb.deleteObjectInput)
}
