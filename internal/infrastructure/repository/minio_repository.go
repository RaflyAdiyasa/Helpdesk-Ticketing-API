package repository

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

type MinioRepository struct {
	internalClient *minio.Client
	publicClient   *minio.Client
	bucket         string
}

func NewMinioRepository(internalClient *minio.Client, publicClient *minio.Client, bucket string) *MinioRepository {
	return &MinioRepository{
		internalClient: internalClient,
		publicClient:   publicClient,
		bucket:         bucket,
	}
}

func (r *MinioRepository) Upload(file multipart.File, objectName string, size int64, contentType string) (string, error) {
	if r.internalClient == nil {
		return "", fmt.Errorf("minio client is not initialized")
	}

	if r.bucket == "" {
		return "", fmt.Errorf("minio bucket is not configured")
	}

	ctx := context.Background()

	exists, err := r.internalClient.BucketExists(ctx, r.bucket)
	if err != nil {
		return "", err
	}

	if !exists {
		if err := r.internalClient.MakeBucket(ctx, r.bucket, minio.MakeBucketOptions{}); err != nil {
			return "", err
		}
	}

	info, err := r.internalClient.PutObject(
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

func (r *MinioRepository) GetPresignedURL(ctx context.Context, objectName string) (string, error) {
	expiry := 15 * time.Minute

	reqParams := make(url.Values)

	presignedURL, err := r.publicClient.PresignedGetObject(
		ctx,
		r.bucket,
		objectName,
		expiry,
		reqParams,
	)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

func (r *MinioRepository) Ping(ctx context.Context) error {
	_, err := r.internalClient.BucketExists(ctx, r.bucket)
	return err
}
