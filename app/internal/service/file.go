package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
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

	tmp, err := os.CreateTemp("", "upload-*")
	if err != nil {
		return nil, err
	}
	defer func(){
		tmp.Close()
		os.Remove(tmp.Name())
	}()

	hasher := sha256.New()
	body := io.Reader(in.Reader)
	_, err = io.Copy(io.MultiWriter(tmp, hasher), body)
	if err != nil {
		return nil, err
	}
	sum := hex.EncodeToString(hasher.Sum(nil))

	if existing, err := s.repo.Files().BySHA256(ctx, sum); err == nil {
		reqlog.From(ctx).Debug("file exists", "file", existing)
		return existing, nil
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	key := shardKey(id.String())

 	if _, err := tmp.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	if err := s.storage.Put(ctx, key, tmp, in.Size, ct); err != nil {
		return nil, err
	}

	f := &domain.File{
		ID: id.String(),
		StorageKey: key,
		OriginalName: name,
		ContentType: ct,
		SHA256: sum,
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

func (s *FileService) ByIDs(ctx context.Context, ids []string) (map[string]*domain.File, error) {
	return s.repo.Files().ByIDs(ctx, ids)
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
