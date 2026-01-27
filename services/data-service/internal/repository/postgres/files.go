package postgres

import (
	"context"
	"fmt"

	"data-service/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FileRepository struct {
	db *pgxpool.Pool
}

func NewFileRepository(db *pgxpool.Pool) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) Create(ctx context.Context, file *models.File) error {
	query := `
INSERT INTO data.files (uuid, filename, size, content_type, bucket, path)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at
	`

	return r.db.QueryRow(ctx, query,
		file.UUID,
		file.Filename,
		file.Size,
		file.ContentType,
		file.Bucket,
		file.Path,
	).Scan(
		&file.ID,
		&file.CreatedAt,
	)
}

func (r *FileRepository) GetByID(ctx context.Context, id int) (*models.File, error) {
	query := `
SELECT id, uuid, filename, size, content_type, bucket, path, created_at
FROM data.files
WHERE id = $1
	`

	file := &models.File{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&file.ID,
		&file.UUID,
		&file.Filename,
		&file.Size,
		&file.ContentType,
		&file.Bucket,
		&file.Path,
		&file.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get file: %v", err)
	}

	return file, nil
}

func (r *FileRepository) GetByUUID(ctx context.Context, uuid string) (*models.File, error) {
	query := `
SELECT id, uuid, filename, size, content_type, bucket, path, created_at
FROM data.files
WHERE uuid = $1
	`

	file := &models.File{}

	err := r.db.QueryRow(ctx, query, uuid).Scan(
		&file.ID,
		&file.UUID,
		&file.Filename,
		&file.Size,
		&file.ContentType,
		&file.Bucket,
		&file.Path,
		&file.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get file: %v", err)
	}

	return file, nil
}

func (r *FileRepository) Delete(ctx context.Context, id int) error {
	query := `
DELETE FROM data.files
WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *FileRepository) DeleteByUUID(ctx context.Context, uuid string) error {
	query := `
DELETE FROM data.files
WHERE uuid = $1
	`
	_, err := r.db.Exec(ctx, query, uuid)
	return err
}
