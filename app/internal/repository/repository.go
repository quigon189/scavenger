package repository

import (
	"context"
	"scavenger/internal/domain"
)

type Repository interface {
	Users() UserRepo
	Sessions() SessionRepo

	WithTx(ctx context.Context, fn func(Repository) error) error

	Close() error
}

type UserRepo interface {
	Create(ctx context.Context, u *domain.User) error
	ByID(ctx context.Context, id int64) (*domain.User, error)
	ByEmail(ctx context.Context, email string) (*domain.User, error)
	List(ctx context.Context) ([]domain.User, error)
}

type SessionRepo interface {
	Create(ctx context.Context, s *domain.Session) error
	ByID(ctx context.Context, id string) (*domain.Session, error)
	Delete(ctx context.Context, id string) error
	DeleteExpired(ctx context.Context) error
}
