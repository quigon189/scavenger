package service

import (
	"context"
	"errors"
	"scavenger/internal/domain"
	"scavenger/internal/repository"

	"github.com/google/uuid"
)

type MaterialService struct {
	repo repository.Repository
}

func NewMaterial(repo repository.Repository) *MaterialService {
	return &MaterialService{repo: repo}
}

type MaterialNoteInput struct {
	DisciplineID int64
	Title        string
	Visible      bool
	DisplayOrder int
	Description  string
}

func (s *MaterialService) CreateNote(ctx context.Context, in MaterialNoteInput) (domain.Material, error) {
	disc, err := s.repo.Disciplines().ByID(ctx, in.DisciplineID)
	if err != nil {
		return nil, err
	}
	if in.Title == "" || in.Description == "" {
		return nil, ErrInvalidTitle
	}

	id, err := s.generateID(ctx)
	if err != nil {
		return nil, err
	}
	
	m := &domain.NoteMaterial{
		BaseMaterial: domain.BaseMaterial{
			ID: id,
			DisciplineID: disc.ID,
			Title: in.Title,
			Visible: in.Visible,
		},
		Description: in.Description,
	}

	if err := s.repo.Materials().Save(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *MaterialService) ListByDisciplineID(ctx context.Context, id int64) ([]domain.Material, error) {
	return s.repo.Materials().ByDisciplineID(ctx, id)
}

func (s *MaterialService) generateID(ctx context.Context) (string,error) {
	for range 10 {
		id, err := uuid.NewV7()
		if err != nil {
			return "", err
		}
		_, err = s.repo.Materials().ByID(ctx, id.String())
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return id.String(), nil
			}
		}
	}
	return "", errors.New("failed to generate id")
}
