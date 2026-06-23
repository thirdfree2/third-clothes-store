package ports

import (
	"context"
	"io"
)

type UploadObjectInput struct {
	Bucket      string
	ObjectKey   string
	Reader      io.Reader
	Size        int64
	ContentType string
}

type ObjectStorage interface {
	Upload(ctx context.Context, input UploadObjectInput) error
}
