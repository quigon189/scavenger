package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"scavenger/internal/domain"
	"scavenger/internal/repository"
	"scavenger/internal/reqlog"
	"scavenger/internal/storage"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrFileTooLarge     = errors.New("file too large")
	ErrInvalideFileName = errors.New("invalid file name")
	ErrEmptyFile        = errors.New("empty file")
)

type FileService struct {
	repo    repository.Repository
	storage storage.Storage
	maxSize int64
}

func NewFile(repo repository.Repository, st storage.Storage, maxSize int64) *FileService {
	return &FileService{
		repo:    repo,
		storage: st,
		maxSize: maxSize,
	}
}

type FileUploadInput struct {
	OwnerID      int64
	OriginalName string
	ContentType  string
	Size         int64
	Reader       io.Reader
}

func (s *FileService) Upload(ctx context.Context, in FileUploadInput) (*domain.File, error) {
	if in.Size <= 0 {
		return nil, ErrEmptyFile
	}
	if s.maxSize > 0 && in.Size > s.maxSize {
		return nil, ErrFileTooLarge
	}
	name := sanitizeName(in.OriginalName)
	if name == "" {
		return nil, ErrInvalideFileName
	}
	ct := in.ContentType
	if ct == "" {
		ct = "application/octet-stream"
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	key := shardKey(id.String())

	if err := s.storage.Put(ctx, key, in.Reader, in.Size, ct); err != nil {
		return nil, err
	}

	f := &domain.File{
		ID: id.String(),
		StorageKey: key,
		OriginalName: name,
		ContentType: ct,
		Size: in.Size,
		OwnerID: in.OwnerID,
	}

	if err := s.repo.Files().Save(ctx, f); err != nil {
		s.storage.Delete(ctx, key)
		return nil, err
	}

	return f, nil
}

func (s *FileService) Open(ctx context.Context, id string) (*domain.File, io.ReadSeekCloser, error) {
	f, err := s.repo.Files().ByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	rc, err := s.storage.Get(ctx, f.StorageKey)
	if err != nil {
		return nil, nil, err
	}
	return f, rc, nil
}

func (s *FileService) ByID(ctx context.Context, id string) (*domain.File, error) {
	return s.repo.Files().ByID(ctx, id)
}

func (s *FileService) Delete(ctx context.Context, id string) error {
	f, err := s.repo.Files().ByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.Files().Delete(ctx, id); err != nil {
		return err
	}
	if err := s.storage.Delete(ctx, f.StorageKey); err != nil {
		log := reqlog.From(ctx)
		log.Warn("delete file from storage", "err", err)
	}
	return nil
}

func sanitizeName(name string) string {
	name = strings.ReplaceAll(name, "\\", "/")
	name = filepath.Base(name)
	name = strings.TrimSpace(name)
	name = strings.Trim(name, ".")
	if name == "/" {
		return ""
	}
	if len(name) > 255 {
		name = name[:255]
	}
	return name
}

func shardKey(id string) string {
	if len(id) < 4 {
		return "files/" + id
	}
	return fmt.Sprintf("files/%s/%s/%s", id[0:2], id[2:4], id)
}
