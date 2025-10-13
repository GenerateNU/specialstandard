package s3_client

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (c *Client) GeneratePresignedURL(ctx context.Context, key string,
	expiry time.Duration) (string, error) {
	psClient := s3.NewPresignClient(c.S3)

	req := &s3.GetObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(key),
	}

	presigned, err := psClient.PresignGetObject(ctx, req, func(opts *s3.PresignOptions) {
		opts.Expires = expiry
	})
	if err != nil {
		return "", fmt.Errorf("failed to presign URL for key %q: %w", key, err)
	}

	return presigned.URL, nil
}
