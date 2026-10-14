package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"scavenger/internal/domain"
	"scavenger/internal/repository"
	"time"
)

type fileRepo struct{ q queryer }

func (r *fileRepo) Save(ctx context.Context, f *domain.File) error {
	creatdAt := time.Now()
	_, err := r.q.ExecContext(ctx,
		`INSERT INTO files (id, storage_key, original_name, content_type, sha256, size, owner_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT (id) DO UPDATE SET
			original_name = excluded.original_name,
			owner_id = excluded.owner_id`,
		f.ID, f.StorageKey, f.OriginalName, f.ContentType, f.SHA256, f.Size, f.OwnerID, creatdAt,
	)
	if err != nil {
		return err
	}
	f.CreatedAt = creatdAt
	return nil
}

func (r *fileRepo) BySHA256(ctx context.Context, sha256 string) (*domain.File, error) {
	return scanFile(r.q.QueryRowContext(ctx,
		`SELECT id, storage_key, original_name, content_type, sha256, size, owner_id, created_at
		FROM files
		WHERE sha256 = ?`, sha256,
	))
}

func (r *fileRepo) ByID(ctx context.Context, id string) (*domain.File, error) {
	return scanFile(r.q.QueryRowContext(ctx,
		`SELECT id, storage_key, original_name, content_type, sha256, size, owner_id, created_at
		FROM files
		WHERE id = ?`, id,
	))
}

func (r *fileRepo) ByIDs(ctx context.Context, ids []string) (map[string]*domain.File, error) {
	if len(ids) == 0 {
		return map[string]*domain.File{}, nil
	}
	q, args := inQuery(
		`SELECT id, storage_key, original_name, content_type, size, size, owner_id, created_at
		FROM files WHERE id IN (%s)`, ids,
	)

	rows, err := r.q.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[string]*domain.File, len(ids))
	for rows.Next() {
		f, err := scanFile(rows)
		if err != nil {
			return nil, err
		}
		out[f.ID] = f
	}
	return out, nil
}

func (r *fileRepo) Delete(ctx context.Context, id string) error {
	_, err := r.q.ExecContext(ctx,
		`DELETE FROM files WHERE id = ?`, id,
	)
	return err
}

func scanFile(row rowScanner) (*domain.File, error) {
	var f domain.File
	err := row.Scan(&f.ID, &f.StorageKey, &f.OriginalName, &f.ContentType, &f.SHA256, &f.Size, &f.OwnerID, &f.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &f, nil
}
