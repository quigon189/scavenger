package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"scavenger/internal/domain"
	"scavenger/internal/repository"
	"time"
)

type materialRepo struct{ q queryer }

type payloadNote struct {
	Description string `json:"description"`
}

type payloadTheory struct {
	Content string `json:"content"`
}

type payloadPractical struct {
	Number   int       `json:"number"`
	Content  string    `json:"content"`
	Deadline time.Time `json:"deadline"`
}

func (r *materialRepo) Save(ctx context.Context, m domain.Material) error {
	var payload any
	switch v := m.(type) {
	case *domain.NoteMaterial:
		payload = payloadNote{Description: v.Description}
	case *domain.TheoryMaterial:
		payload = payloadTheory{Content: v.Content}
	case *domain.PracticalMaterial:
		payload = payloadPractical{
			Number:   v.Number,
			Content:  v.Content,
			Deadline: v.Deadline,
		}
	default:
		return repository.ErrUnknownMaterialType
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = r.q.ExecContext(ctx,
		`INSERT INTO materials (id, discipline_id, title, display_order, visible, type, payload, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			title         = excluded.title,
			discipline_id = excluded.discipline_id,
			display_order = excluded.display_order,
			visible       = excluded.visible,
			payload       = excluded.payload,
			updated_at    = excluded.created_at`,
		m.GetID(),
		m.GetDisciplineID(),
		m.GetTitle(),
		m.GetDisplayOrder(),
		m.GetVisible(),
		m.Kind(),
		string(payloadJSON),
		time.Now().UTC())
	return err
}

func (r *materialRepo) ByID(ctx context.Context, id string) (domain.Material, error) {
	return scanMaterial(r.q.QueryRowContext(ctx,
		`SELECT id, title, discipline_id, display_order, visible, type, payload, created_at, updated_at
		FROM materials
		WHERE id = ?`, id))
}

func (r *materialRepo) ByDisciplineID(ctx context.Context, discID int64) ([]domain.Material, error) {
	rows, err := r.q.QueryContext(ctx,
		`SELECT id, title, discipline_id, display_order, visible, type, payload, created_at, updated_at
		FROM materials
		WHERE discipline_id = ?`, discID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materials []domain.Material
	for rows.Next() {
		m, err := scanMaterial(rows)
		if err != nil {
			return nil, err
		}
		materials = append(materials, m)
	}
	return materials, nil
}

func (r *materialRepo) Delete(ctx context.Context, id string) error {
	_, err := r.q.ExecContext(ctx, `DELETE FROM materials WHERE id = ?`, id)
	return err
}

func scanMaterial(s rowScanner) (domain.Material, error) {
	var (
		id, title, typ, payloadRaw string
		discID                     int64
		dispOrder                  int
		visible                    bool
		createdAt                  time.Time
		updatedAt                  *time.Time
	)

	err := s.Scan(&id, &title, &discID, &dispOrder, &visible, &typ,
		&payloadRaw, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	baseMaterial := domain.BaseMaterial{
		ID:           id,
		DisciplineID: discID,
		DisplayOrder: dispOrder,
		Title:        title,
		Visible:      visible,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}

	switch domain.MaterialKind(typ) {
	case domain.MaterialNoteType:
		var p payloadNote
		if err := json.Unmarshal([]byte(payloadRaw), &p); err != nil {
			return nil, err
		}
		return domain.NoteMaterial{
			BaseMaterial: baseMaterial,
			Description:  p.Description,
		}, nil
	case domain.MaterialTheoryType:
		var p payloadTheory
		if err := json.Unmarshal([]byte(payloadRaw), &p); err != nil {
			return nil, err
		}
		return domain.TheoryMaterial{
			BaseMaterial: baseMaterial,
			Content: p.Content,
		}, nil
	case domain.MaterialPracticalType:
		var p payloadPractical
		if err := json.Unmarshal([]byte(payloadRaw), &p); err != nil {
			return nil, err
		}
		return domain.PracticalMaterial{
			BaseMaterial: baseMaterial,
			Number: p.Number,
			Content: p.Content,
			Deadline: p.Deadline,
		}, nil
	}

	return nil, repository.ErrUnknownMaterialType
}
