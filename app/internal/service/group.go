package service

import (
	"context"
	"errors"
	"scavenger/internal/domain"
	"scavenger/internal/repository"
)

var (
	ErrInvalidInput = errors.New("invalid group input")
)

type GroupService struct {
	repo repository.Repository
}

func NewGroup(repo repository.Repository) *GroupService {
	return &GroupService{repo: repo}
}

type GroupInput struct {
	Number          int
	StartYear       int
	DurationOfStudy int
	Specialty       string
	ShortSpecialty  string
}

func (s *GroupService) Create(ctx context.Context, in GroupInput) (*domain.Group, error) {
	if in.Number == 0 {
		return nil, ErrInvalidInput
	}
	if in.StartYear < 2000 || in.StartYear > 3000 {
		return nil, ErrInvalidInput
	}
	if in.DurationOfStudy == 0 {
		return nil, ErrInvalidInput
	}
	if in.Specialty == "" || in.ShortSpecialty == "" {
		return nil, ErrInvalidInput
	}

	group := domain.Group{
		Number:         int64(in.Number),
		StartYear:      int64(in.StartYear),
		EndYear:        int64(in.StartYear + in.DurationOfStudy),
		Specialty:      in.Specialty,
		ShortSpecialty: in.ShortSpecialty,
	}

	if err := s.repo.Groups().Create(ctx, &group); err != nil {
		return nil, err
	}
	return &group, nil
}

func (s *GroupService) List(ctx context.Context) ([]domain.Group, error) {
	return s.repo.Groups().List(ctx)
}

func (s *GroupService) Students(ctx context.Context, gid int64) ([]domain.User, error) {
	var students []domain.User

	err := s.repo.WithTx(ctx, func(r repository.Repository) error {
		studentIDs, err := r.Groups().StudentIDs(ctx, gid)
		if err != nil {
			return err
		}

		for _, id := range studentIDs {
			student, err := r.Users().ByID(ctx, id)
			if err != nil {
				return err
			}

			students = append(students, *student)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return students, nil
}

func (s *GroupService) AddStudent(ctx context.Context, gid int64, sid int64) error {
	return s.repo.Groups().AddStudent(ctx, gid, sid)
}
