package fs

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"scavenger/internal/storage"
	"strings"
)

type Storage struct {
	root string
}

func New(root string) (*Storage, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(abs, 0o755); err != nil {
		return nil, err
	}
	return &Storage{root: abs}, nil
}

func (s Storage) Put(ctx context.Context, key string, r io.Reader, _ int64, _ string) error {
	full, err := s.resolve(key)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return err
	}

	tmp, err := os.CreateTemp(filepath.Dir(full), ".upload-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() {
		tmp.Close()
		os.Remove(tmpName)
	}()

	if _, err := io.Copy(tmp, r); err != nil {
		return err
	}
	if err := tmp.Sync(); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, full)
}

func (s *Storage) Get(ctx context.Context, key string) (io.ReadSeekCloser, error) {
	full, err := s.resolve(key)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(full)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}
	return f, nil
}

func (s *Storage) Delete(ctx context.Context, key string) error {
	full, err := s.resolve(key)
	if err != nil {
		return err
	}
	if err := os.Remove(full); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return storage.ErrNotFound
		}
		return err
	}
	return nil
}

func (s *Storage) Stat(ctx context.Context, key string) (*storage.ObjectInfo, error ) {
	full, err := s.resolve(key)
	if err != nil {
		return nil, err
	}
	fi, err := os.Stat(full)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}
	return &storage.ObjectInfo{
		Size: fi.Size(),
		ModTime: fi.ModTime(),
	}, nil
}

func (s Storage) resolve(key string) (string, error) {
	if key == "" || strings.HasPrefix(key, "/") {
		return "", storage.ErrInvalidKey
	}
	clean := filepath.Clean(key)
	if clean == "." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) || clean == ".." {
		return "", storage.ErrInvalidKey
	}
	full := filepath.Join(s.root, clean)
	return full, nil
}
