package minioadapter

import (
	"context"

	"clothing-service/internal/application/ports"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ObjectStorage struct {
	client *minio.Client
}

func NewObjectStorage(endpoint, accessKey, secretKey string, useSSL bool) (*ObjectStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return &ObjectStorage{
		client: client,
	}, nil
}

func (s *ObjectStorage) EnsureBucket(ctx context.Context, bucket string) error {
	exists, err := s.client.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	return s.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
}

func (s *ObjectStorage) Upload(ctx context.Context, input ports.UploadObjectInput) error {
	_, err := s.client.PutObject(
		ctx,
		input.Bucket,
		input.ObjectKey,
		input.Reader,
		input.Size,
		minio.PutObjectOptions{
			ContentType: input.ContentType,
		},
	)

	return err
}
