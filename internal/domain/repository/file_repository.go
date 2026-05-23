package repository

import (
	"context"
	"mime/multipart"
)

type FileRepository interface {
	Upload(file multipart.File, objectName string, size int64, contentType string) (string, error)
	GetPresignedURL(ctx context.Context, objectName string) (string, error)
}
