package repository

import (
	"context"
	"scavenger/internal/domain"
)

type Repository interface {
	Users() UserRepo
	Sessions() SessionRepo

	Groups() GroupRepo
	Disciplines() DisciplineRepo
	Materials() MaterialRepo

	WithTx(ctx context.Context, fn func(Repository) error) error

	Close() error
}

type UserRepo interface {
	Save(ctx context.Context, u *domain.User) error
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

type GroupRepo interface {
	Save(ctx context.Context, g *domain.Group) error
	ByID(ctx context.Context, id int64) (*domain.Group, error)
	ByStudentID(ctx context.Context, id int64) (*domain.Group, error)
	List(ctx context.Context) ([]domain.Group, error)
	StudentIDs(ctx context.Context, id int64) ([]int64, error)
	AddStudent(ctx context.Context, gid int64, sid int64) error
	RemoveStudent(ctx context.Context, sid int64) error
}

type DisciplineRepo interface {
	Save(ctx context.Context, d *domain.Discipline) error
	ByID(ctx context.Context, id int64) (*domain.Discipline, error)
	ByTeacherID(ctx context.Context, id int64) ([]domain.Discipline, error)
	ByGroupID(ctx context.Context, id int64) ([]domain.Discipline, error)
	List(ctx context.Context) ([]domain.Discipline, error)
	Delete(ctx context.Context, id int64) error
}

type MaterialRepo interface {
	Save(ctx context.Context, m domain.Material) error
	ByID(ctx context.Context, id string) (domain.Material, error)
	ByDisciplineID(ctx context.Context, id int64) ([]domain.Material, error)
	Delete(ctx context.Context, id string) error
}
