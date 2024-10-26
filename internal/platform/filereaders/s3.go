package filereaders

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3FileReader struct {
	S3Client *s3.Client
	Bucket   string
}

func NewS3FileReader(client *s3.Client, bucket string) *S3FileReader {
	return &S3FileReader{
		S3Client: client,
		Bucket:   bucket,
	}
}

func (s *S3FileReader) Open(ctx context.Context, filePath string) (io.ReadCloser, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(filePath),
	}

	output, err := s.S3Client.GetObject(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return output.Body, nil
}
