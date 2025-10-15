package s3_client

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

func (c *Client) ListObjectsByPrefix(ctx context.Context, prefix string) ([]string, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(c.Bucket),
		Prefix: aws.String(prefix),
	}

	res, err := c.S3.ListObjectsV2(ctx, input)
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.ErrorCode() {
			case "AccessDenied":
				return nil, fmt.Errorf("access denied to bucket %q: %w", prefix, err)
			case "NoSuchBucket":
				return nil, fmt.Errorf("bucket %q does not exist: %w", c.Bucket, err)
			default:
				return nil, fmt.Errorf("failed to list objects for prefix %q: %w", prefix, err)
			}
		}
	}

	var keys []string
	for _, obj := range res.Contents {
		keys = append(keys, aws.ToString(obj.Key))
	}

	return keys, nil
}
