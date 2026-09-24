package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"scavenger/internal/domain"
	"scavenger/internal/repository"
	"time"
)

type sessionRepo struct{ q queryer }

func (r *sessionRepo) Create(ctx context.Context, s *domain.Session) error {
	_, err := r.q.ExecContext(ctx,
		`INSERT INTO sessions(id, user_id, expires_at, created_at)
		VALUES (?, ?, ?, ?)`,
		s.ID, s.UserID, s.ExpiresAt, time.Now().UTC(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *sessionRepo) ByID(ctx context.Context, id string) (*domain.Session, error) {
	row := r.q.QueryRowContext(ctx,
		`SELECT id, user_id, expires_at, created_at
		FROM sessions WHERE id = ?`, id,
	)
	var s domain.Session
	if err := row.Scan(&s.ID, &s.UserID, &s.ExpiresAt, &s.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *sessionRepo) Delete(ctx context.Context, id string) error {
	_, err := r.q.ExecContext(ctx,
		`DELETE FROM sessions WHERE id = ?`, id,
	)
	return err
}

func (r *sessionRepo) DeleteExpired(ctx context.Context) error {
	_, err := r.q.ExecContext(ctx,
		`DELETE FROM sessions WHERE expires_at < datetime('now')`,
	)
	return err
}
