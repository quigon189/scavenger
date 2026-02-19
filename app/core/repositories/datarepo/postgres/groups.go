package postgres

import (
	"context"
	"fmt"

	"scavenger/core/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GroupRepository struct {
	db *pgxpool.Pool
}

func NewGroupRepository(db *pgxpool.Pool) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) Create(ctx context.Context, group *models.Group) error {
	query := `
INSERT INTO data.groups (name)
VALUES ($1)
RETURNING id, created_at
	`

	return r.db.QueryRow(ctx, query, group.Name).Scan(
		&group.ID,
		&group.CreatedAt,
	)
}

func (r *GroupRepository) GetByID(ctx context.Context, id int) (*models.Group, error) {
	query := `
SELECT id, name, created_at
FROM data.groups
WHERE id = $1
	`

	group := &models.Group{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&group.ID,
		&group.Name,
		&group.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get group: %v", err)
	}

	return group, nil
}

func (r *GroupRepository) GetByName(ctx context.Context, name string) (*models.Group, error) {
	query := `
SELECT id, name, created_at
FROM data.groups
WHERE name = $1
	`

	group := &models.Group{}

	err := r.db.QueryRow(ctx, query, name).Scan(
		&group.ID,
		&group.Name,
		&group.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get group: %v", err)
	}

	return group, nil
}

func (r *GroupRepository) GetWithStudents(ctx context.Context, id int) (*models.Group, error) {
	group, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	sr := NewStudentRepository(r.db)
	students, err := sr.GetByGroupID(ctx, id)
	if err != nil {
		return nil, err
	}
	group.Students = append(group.Students, students...)

	return group, nil
}

func (r *GroupRepository) Update(ctx context.Context, group *models.Group) error {
	query := `
UPDATE data.groups
SET name = $1
WHERE id = $2
	`

	_, err := r.db.Exec(ctx, query, group.Name, group.ID)
	return err
}

func (r *GroupRepository) Delete(ctx context.Context, id int) error {
	query := `
DELETE FROM data.groups
WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *GroupRepository) GetAll(ctx context.Context) ([]models.Group, error) {
	query := `
SELECT id, name, created_at
FROM data.groups
ORDER BY name
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query groups: %v", err)
	}
	defer rows.Close()

	var groups []models.Group

	for rows.Next() {
		var group models.Group

		err := rows.Scan(
			&group.ID,
			&group.Name,
			&group.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan group: %v", err)
		}

		groups = append(groups, group)
	}

	return groups, nil
}
