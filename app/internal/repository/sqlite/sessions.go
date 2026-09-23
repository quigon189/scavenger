package sqlite

import (
	"context"
	"scavenger/internal/domain"
)

type sessionRepo struct {q queryer}

func (r *sessionRepo) Create(ctx context.Context, s *domain.Session) error {
	return nil
}
