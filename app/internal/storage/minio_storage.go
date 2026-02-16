package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type MinioStorage struct {
	client     *minio.Client
	bucketName string
}

func NewMinioStorage(client *minio.Client, bucketName string) (*MinioStorage, error) {
	return &MinioStorage{
		client:     client,
		bucketName: bucketName,
	}, nil
}

func (s *MinioStorage) Upload(ctx context.Context, filename, fileType, contentType string, reader io.Reader, size int64) (*FileInfo, error) {
	fileUUID := uuid.New().String()
	
	path := fmt.Sprintf("%s/%s", fileType, fileUUID)
	
	_, err := s.client.PutObject(ctx, s.bucketName, path, reader, size, minio.PutObjectOptions{
		ContentType:        contentType,
		ContentDisposition: fmt.Sprintf("attachment; filename=\"%s\"", filename),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to MinIO: %w", err)
	}
	
	return &FileInfo{
		UUID:        fileUUID,
		Filename:    filename,
		Size:        size,
		ContentType: contentType,
		Bucket:      s.bucketName,
		Path:        path,
		CreatedAt:   time.Now(),
	}, nil
}

func (s *MinioStorage) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	object, err := s.client.GetObject(ctx, s.bucketName, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from MinIO: %w", err)
	}
	
	_, err = object.Stat()
	if err != nil {
		object.Close()
		return nil, fmt.Errorf("object not found: %w", err)
	}
	
	return object, nil
}

func (s *MinioStorage) GetSignedURL(ctx context.Context, path string, expires time.Duration) (string, error) {
	if expires <= 0 {
		expires = 15 * time.Minute // дефолтное время жизни ссылки
	}
	
	url, err := s.client.PresignedGetObject(ctx, s.bucketName, path, expires, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}
	
	return url.String(), nil
}

func (s *MinioStorage) Delete(ctx context.Context, path string) error {
	err := s.client.RemoveObject(ctx, s.bucketName, path, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object from MinIO: %w", err)
	}
	
	return nil
}

func (s *MinioStorage) Stat(ctx context.Context, path string) (*FileInfo, error) {
	objInfo, err := s.client.StatObject(ctx, s.bucketName, path, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to stat object: %w", err)
	}
	
	filename := path
	if objInfo.Metadata.Get("X-Amz-Meta-Filename") != "" {
		filename = objInfo.Metadata.Get("X-Amz-Meta-Filename")
	}
	
	return &FileInfo{
		Filename:    filename,
		Size:        objInfo.Size,
		ContentType: objInfo.ContentType,
		Bucket:      s.bucketName,
		Path:        path,
		CreatedAt:   objInfo.LastModified,
	}, nil
}
