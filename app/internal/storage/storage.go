package storage

import (
	"context"
	"errors"
	"io"
	"time"
)

var (
	ErrNotFound   = errors.New("object not found")
	ErrInvalidKey = errors.New("invalid key")
)

type ObjectInfo struct {
	Size        int64
	ContentType string
	ModTime     time.Time
}

type Storage interface {
	Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error
	Get(ctx context.Context, key string) (io.ReadSeekCloser, error)
	Delete(ctx context.Context, key string) error
	Stat(ctx context.Context, key string) (*ObjectInfo, error)
}
