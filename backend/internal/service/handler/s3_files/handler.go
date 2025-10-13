package s3_files

import "specialstandard/internal/s3_client"

type Handler struct {
	S3Client *s3_client.Client
}

func NewHandler(s3Client *s3_client.Client) *Handler {
	return &Handler{
		S3Client: s3Client,
	}
}
