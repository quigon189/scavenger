package service

import (
	"context"
	"errors"
	"scavenger/internal/domain"
	"scavenger/internal/repository"
)

var (
	ErrInvalidTitle = errors.New("invalid discipline title")
)

type DisciplineService struct {
	repo repository.Repository
} 

func NewDiscipline(repo repository.Repository) *DisciplineService {
	return &DisciplineService{repo: repo}
}

type DisciplineInput struct {
	Title string
	Description string
	TeacherID int64
	GroupID int64
	Archived bool
}

func (s *DisciplineService) Create(ctx context.Context, in DisciplineInput) (*domain.Discipline, error) {
	if in.Title == "" {
		return nil, ErrInvalidTitle
	}
	group, err := s.repo.Groups().ByID(ctx, in.GroupID)
	if err != nil {
		return nil, err
	}
	teacher, err := s.repo.Users().ByID(ctx, in.TeacherID)
	if err != nil {
		return nil, err
	}
	discipline := domain.Discipline{
		Title: in.Title,
		Description: in.Description,
		TeacherID: teacher.ID,
		GroupID: group.ID,
		Archived: in.Archived,
	}

	if err := s.repo.Disciplines().Save(ctx, &discipline); err != nil {
		return nil, err
	}

	return &discipline, nil
}

func (s *DisciplineService) ByID(ctx context.Context, id int64) (*domain.Discipline, error) {
	return s.repo.Disciplines().ByID(ctx, id)
}

func (s *DisciplineService) ByTeacherID(ctx context.Context, teacherID int64) ([]domain.Discipline, error) {
	return s.repo.Disciplines().ByTeacherID(ctx, teacherID)
}

func (s *DisciplineService) ByGroupID(ctx context.Context, groupID int64) ([]domain.Discipline, error) {
	return s.repo.Disciplines().ByGroupID(ctx, groupID)
}

func (s *DisciplineService) List(ctx context.Context) ([]domain.Discipline, error) {
	return s.repo.Disciplines().List(ctx)
}

func (s *DisciplineService) Update(ctx context.Context, d *domain.Discipline) error {
	return s.repo.Disciplines().Save(ctx, d)
}
