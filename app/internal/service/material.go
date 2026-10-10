package service

import (
	"context"
	"errors"
	"scavenger/internal/domain"
	"scavenger/internal/repository"
	"scavenger/internal/reqlog"
	"slices"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidInput = errors.New("invalid input")
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
		return nil, ErrInvalidInput
	}

	id, err := s.generateID(ctx)
	if err != nil {
		return nil, err
	}

	m := &domain.NoteMaterial{
		BaseMaterial: domain.BaseMaterial{
			ID:           id,
			DisciplineID: disc.ID,
			Title:        in.Title,
			Visible:      in.Visible,
			DisplayOrder: in.DisplayOrder,
		},
		Description: in.Description,
	}

	if err := s.repo.Materials().Save(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

type MaterialTheoryInput struct {
	DisciplineID int64
	Title        string
	Visible      bool
	DisplayOrder int
	Content      string
}

func (s *MaterialService) CreateTheory(ctx context.Context, in MaterialTheoryInput) (domain.Material, error) {
	disc, err := s.repo.Disciplines().ByID(ctx, in.DisciplineID)
	if err != nil {
		return nil, err
	}

	if in.Title == "" || in.Content == "" {
		return nil, ErrInvalidInput
	}

	id, err := s.generateID(ctx)
	if err != nil {
		return nil, err
	}

	m := &domain.TheoryMaterial{
		BaseMaterial: domain.BaseMaterial{
			ID:           id,
			DisciplineID: disc.ID,
			Title:        in.Title,
			Visible:      in.Visible,
			DisplayOrder: in.DisplayOrder,
		},
		Content: in.Content,
	}

	if err := s.repo.Materials().Save(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

type MaterialPracticalInput struct {
	DisciplineID int64
	Title        string
	Visible      bool
	DisplayOrder int
	Number       int
	Content      string
	Deadline     time.Time
}

func (s *MaterialService) CreatePractical(ctx context.Context, in MaterialPracticalInput) (domain.Material, error) {
	disc, err := s.repo.Disciplines().ByID(ctx, in.DisciplineID)
	if err != nil {
		return nil, err
	}

	if in.Title == "" || in.Number <= 0 || in.Content == "" || in.Deadline.Before(time.Now()) {
		return nil, ErrInvalidInput
	}

	id, err := s.generateID(ctx)
	if err != nil {
		return nil, err
	}

	m := &domain.PracticalMaterial{
		BaseMaterial: domain.BaseMaterial{
			ID:           id,
			DisciplineID: disc.ID,
			Title:        in.Title,
			Visible:      in.Visible,
			DisplayOrder: in.DisplayOrder,
		},
		Number:   in.Number,
		Content:  in.Content,
		Deadline: in.Deadline,
	}

	if err := s.repo.Materials().Save(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *MaterialService) ListByDisciplineID(ctx context.Context, id int64) ([]domain.Material, error) {
	return s.repo.Materials().ByDisciplineID(ctx, id)
}

func (s *MaterialService) ByID(ctx context.Context, id string) (domain.Material, error) {
	return s.repo.Materials().ByID(ctx, id)
}

func (s *MaterialService) generateID(ctx context.Context) (string, error) {
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

func (s *MaterialService) Move(ctx context.Context, id string, shift int) error {
	log := reqlog.From(ctx)
	return s.repo.WithTx(ctx, func(r repository.Repository) error {
		cur, err := r.Materials().ByID(ctx, id)
		if err != nil {
			return err
		}
		materials, err := r.Materials().ByDisciplineID(ctx, cur.GetDisciplineID())
		if err != nil {
			return err
		}
		slices.SortFunc(materials, func(a, b domain.Material) int {
			return a.GetDisplayOrder() - b.GetDisplayOrder()
		})

		mIdx := []int{}
		for _, m := range materials {
			mIdx = append(mIdx, m.GetDisplayOrder())
		}
		log.Debug("materials", "materials_orders", mIdx)

		curIdx := 0
		for i, m := range materials {
			if m.GetID() == cur.GetID() {
				curIdx = i
			}
		}

		newIdx := curIdx + shift
		if newIdx < 0 {
			newIdx = 0
		}
		if newIdx > len(materials)-1 {
			newIdx = len(materials)-1
		}
		a := materials[0].GetDisplayOrder() - 100
		b := materials[len(materials)-1].GetDisplayOrder() + 100
		switch {
		case newIdx == 0:
			b = materials[0].GetDisplayOrder()
		case newIdx == len(materials)-1:
			a = materials[len(materials)-1].GetDisplayOrder()
		case newIdx < curIdx:
			a = materials[newIdx-1].GetDisplayOrder()
			b = materials[newIdx].GetDisplayOrder()
		default:
			a = materials[newIdx].GetDisplayOrder()
			b = materials[newIdx+1].GetDisplayOrder()
		}
		if b-a == 1 || b == a {
			materials = slices.Delete(materials, curIdx, curIdx+1)
			materials = slices.Insert(materials, newIdx, cur)
			for i, m := range materials {
				m.SetDisplayOrder(100*i + 100)
				if err := r.Materials().Save(ctx, m); err != nil {
					return err
				}
			}
			return nil
		}
		cur.SetDisplayOrder(a + (b-a)/2)
		if err := r.Materials().Save(ctx, cur); err != nil {
			return err
		}
		return nil
	})
}
