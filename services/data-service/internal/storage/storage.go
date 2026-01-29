package storage

import (
	"context"
	"io"
	"time"
)

// FileInfo содержит информацию о файле в хранилище
type FileInfo struct {
	UUID        string    `json:"uuid"`
	Filename    string    `json:"filename"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	Bucket      string    `json:"bucket"`
	Path        string    `json:"path"`
	CreatedAt   time.Time `json:"created_at"`
}

// Storage интерфейс для работы с файловым хранилищем
type Storage interface {
	// Upload загружает файл в хранилище
	Upload(ctx context.Context, filename string, fileType, contentType string, reader io.Reader, size int64) (*FileInfo, error)

	// Download скачивает файл из хранилища
	Download(ctx context.Context, path string) (io.ReadCloser, error)

	// GetSignedURL возвращает подписанную URL для доступа к файлу
	GetSignedURL(ctx context.Context, path string, expires time.Duration) (string, error)

	// Delete удаляет файл из хранилища
	Delete(ctx context.Context, path string) error

	// Stat возвращает информацию о файле
	Stat(ctx context.Context, path string) (*FileInfo, error)
}
