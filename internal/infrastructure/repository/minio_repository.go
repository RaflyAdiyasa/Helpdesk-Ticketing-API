package repository

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
)

type MinioRepository struct {
	client *minio.Client
	bucket string
}

func NewMinioRepository(client *minio.Client, bucket string) *MinioRepository {
	return &MinioRepository{
		client: client,
		bucket: bucket,
	}
}

func (r *MinioRepository) Upload(file multipart.File, objectName string, size int64, contentType string) (string, error) {
	if r.client == nil {
		return "", fmt.Errorf("minio client is not initialized")
	}

	if r.bucket == "" {
		return "", fmt.Errorf("minio bucket is not configured")
	}

	ctx := context.Background()

	exists, err := r.client.BucketExists(ctx, r.bucket)
	if err != nil {
		return "", err
	}

	if !exists {
		if err := r.client.MakeBucket(ctx, r.bucket, minio.MakeBucketOptions{}); err != nil {
			return "", err
		}
	}

	info, err := r.client.PutObject(
		ctx,
		r.bucket,
		objectName,
		file,
		size,
		minio.PutObjectOptions{ContentType: contentType},
	)

	if err != nil {
		return "", err
	}

	return info.Key, nil
}
