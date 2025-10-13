package s3_client

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (c *Client) ListObjectsByPrefix(ctx context.Context, prefix string) ([]string, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(c.Bucket),
		Prefix: aws.String(prefix),
	}

	res, err := c.S3.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects for prefix %q: %w", prefix, err)
	}

	var keys []string
	for _, obj := range res.Contents {
		keys = append(keys, aws.ToString(obj.Key))
	}

	return keys, nil
}
