package services

import (
	"context"
	"data-service/internal/models"
	"data-service/internal/storage"
	"fmt"
	"io"
	"time"
)

type FileService struct {
	service    *Service
	storage    storage.Storage
	signedURLTTL time.Duration
}

func NewFileService(service *Service, storage storage.Storage) *FileService {
	return &FileService{
		service:     service,
		storage:     storage,
		signedURLTTL: 15 * time.Minute, // TTL для подписанных ссылок
	}
}

func (f *FileService) UploadFile(ctx context.Context, filename string, fileType string, contentType string, size int64, reader io.Reader) (*models.File, error) {
	_, err := f.service.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	if fileType == "" {
		return nil, fmt.Errorf("invalid file type")
	}

	// Загружаем файл в хранилище
	fileInfo, err := f.storage.Upload(ctx, filename, fileType, contentType, reader, size)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Создаем запись в базе данных
	file := &models.File{
		UUID:        fileInfo.UUID,
		Filename:    fileInfo.Filename,
		Size:        fileInfo.Size,
		ContentType: fileInfo.ContentType,
		Bucket:      fileInfo.Bucket,
		Path:        fileInfo.Path,
		CreatedAt:   fileInfo.CreatedAt,
	}

	err = f.service.repo.Files.Create(ctx, file)
	if err != nil {
		// Удаляем файл из хранилища в случае ошибки
		f.storage.Delete(ctx, fileInfo.Path)
		return nil, fmt.Errorf("failed to create file record: %w", err)
	}

	return file, nil
}

func (f *FileService) GetFile(ctx context.Context, fileID int) (*models.File, io.ReadCloser, error) {
	_, err := f.service.Authenticate(ctx)
	if err != nil {
		return nil, nil, err
	}

	file, err := f.service.repo.Files.GetByID(ctx, fileID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get file record: %w", err)
	}

	if file == nil {
		return nil, nil, fmt.Errorf("file not found")
	}

	// Скачиваем файл из хранилища
	reader, err := f.storage.Download(ctx, file.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to download file from storage: %w", err)
	}

	return file, reader, nil
}

func (f *FileService) GetFileWithSignedURL(ctx context.Context, fileID int) (*models.File, error) {
	_, err := f.service.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	file, err := f.service.repo.Files.GetByID(ctx, fileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get file record: %w", err)
	}

	if file == nil {
		return nil, fmt.Errorf("file not found")
	}

	// Генерируем подписанную ссылку
	signedURL, err := f.storage.GetSignedURL(ctx, file.Path, f.signedURLTTL)
	if err != nil {
		return file, nil // Возвращаем файл без ссылки, если не удалось сгенерировать
	}

	file.SignedURL = signedURL
	return file, nil
}

func (f *FileService) DeleteFile(ctx context.Context, fileID int) error {
	// Только админ может удалять файлы
	_, err := f.service.RequireRole(ctx, "admin")
	if err != nil {
		return err
	}

	file, err := f.service.repo.Files.GetByID(ctx, fileID)
	if err != nil {
		return fmt.Errorf("failed to get file record: %w", err)
	}

	if file == nil {
		return fmt.Errorf("file not found")
	}

	// Удаляем из хранилища
	err = f.storage.Delete(ctx, file.Path)
	if err != nil {
		return fmt.Errorf("failed to delete file from storage: %w", err)
	}

	// Удаляем из базы данных
	return f.service.repo.Files.Delete(ctx, fileID)
}

func (f *FileService) GenerateSignedURL(ctx context.Context, fileID int, expires time.Duration) (string, error) {
	_, err := f.service.Authenticate(ctx)
	if err != nil {
		return "", err
	}

	file, err := f.service.repo.Files.GetByID(ctx, fileID)
	if err != nil {
		return "", fmt.Errorf("failed to get file record: %w", err)
	}

	if file == nil {
		return "", fmt.Errorf("file not found")
	}

	if expires <= 0 {
		expires = f.signedURLTTL
	}

	return f.storage.GetSignedURL(ctx, file.Path, expires)
}
