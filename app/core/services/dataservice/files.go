package dataservice

import (
	"context"
	"scavenger/core/models"
	"fmt"
	"io"
	"time"
)

func (s *DataService) UploadFile(ctx context.Context, filename string, fileType string, contentType string, size int64, reader io.Reader) (*models.File, error) {
	_, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	if fileType == "" {
		return nil, fmt.Errorf("invalid file type")
	}

	// Загружаем файл в хранилище
	fileInfo, err := s.storage.Upload(ctx, filename, fileType, contentType, reader, size)
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

	err = s.dataRepo.Files.Create(ctx, file)
	if err != nil {
		// Удаляем файл из хранилища в случае ошибки
		s.storage.Delete(ctx, fileInfo.Path)
		return nil, fmt.Errorf("failed to create file record: %w", err)
	}

	return file, nil
}

func (s *DataService) GetFile(ctx context.Context, fileID int) (*models.File, io.ReadCloser, error) {
	_, err := s.Authenticate(ctx)
	if err != nil {
		return nil, nil, err
	}

	file, err := s.dataRepo.Files.GetByID(ctx, fileID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get file record: %w", err)
	}

	if file == nil {
		return nil, nil, fmt.Errorf("file not found")
	}

	// Скачиваем файл из хранилища
	reader, err := s.storage.Download(ctx, file.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to download file from storage: %w", err)
	}

	return file, reader, nil
}

func (s *DataService) GetFileWithSignedURL(ctx context.Context, fileID int) (*models.File, error) {
	_, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	file, err := s.dataRepo.Files.GetByID(ctx, fileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get file record: %w", err)
	}

	if file == nil {
		return nil, fmt.Errorf("file not found")
	}

	// Генерируем подписанную ссылку
	signedURL, err := s.storage.GetSignedURL(ctx, file.Path, s.signedURLTTL)
	if err != nil {
		return file, nil // Возвращаем файл без ссылки, если не удалось сгенерировать
	}

	file.SignedURL = signedURL
	return file, nil
}

func (s *DataService) DeleteFile(ctx context.Context, fileID int) error {
	// Только админ может удалять файлы
	_, err := s.RequireRole(ctx, "admin")
	if err != nil {
		return err
	}

	file, err := s.dataRepo.Files.GetByID(ctx, fileID)
	if err != nil {
		return fmt.Errorf("failed to get file record: %w", err)
	}

	if file == nil {
		return fmt.Errorf("file not found")
	}

	// Удаляем из хранилища
	err = s.storage.Delete(ctx, file.Path)
	if err != nil {
		return fmt.Errorf("failed to delete file from storage: %w", err)
	}

	// Удаляем из базы данных
	return s.dataRepo.Files.Delete(ctx, fileID)
}

func (s *DataService) GenerateSignedURL(ctx context.Context, fileID int, expires time.Duration) (string, error) {
	_, err := s.Authenticate(ctx)
	if err != nil {
		return "", err
	}

	file, err := s.dataRepo.Files.GetByID(ctx, fileID)
	if err != nil {
		return "", fmt.Errorf("failed to get file record: %w", err)
	}

	if file == nil {
		return "", fmt.Errorf("file not found")
	}

	if expires <= 0 {
		expires = s.signedURLTTL
	}

	return s.storage.GetSignedURL(ctx, file.Path, expires)
}
