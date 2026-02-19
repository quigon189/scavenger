package postgres

import (
	"context"
	"scavenger/core/models"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DisciplineRepository struct {
	db *pgxpool.Pool
}

func NewDisciplineRepository(db *pgxpool.Pool) *DisciplineRepository {
	return &DisciplineRepository{db: db}
}

func (r *DisciplineRepository) Create(ctx context.Context, discipline *models.Discipline) error {
	query := `
INSERT INTO data.disciplines (name, teacher_id, group_id, period_id, description)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(ctx, query,
		discipline.Name,
		discipline.TeacherID,
		discipline.GroupID,
		discipline.PeriodID,
		discipline.Description,
	).Scan(
		&discipline.ID,
		&discipline.CreatedAt,
		&discipline.UpdatedAt,
	)
}

func (r *DisciplineRepository) GetByID(ctx context.Context, id int) (*models.Discipline, error) {
	query := `
SELECT d.id, d.name, d.teacher_id, d.group_id, d.period_id, d.description, 
       d.created_at, d.updated_at,
       t.id, t.created_at, t.updated_at,
       u.username, u.name, u.role,
       g.id, g.name, g.created_at,
       p.id, p.name, p.half_year, p.start_date, p.end_date, p.created_at
FROM data.disciplines d
LEFT JOIN data.teachers t ON t.id = d.teacher_id
LEFT JOIN auth.users u ON u.id = t.id
LEFT JOIN data.groups g ON g.id = d.group_id
LEFT JOIN data.periods p ON p.id = d.period_id
WHERE d.id = $1
	`

	discipline := &models.Discipline{}
	var teacherCreatedAt, teacherUpdatedAt time.Time
	var groupCreatedAt time.Time
	var periodCreatedAt time.Time

	err := r.db.QueryRow(ctx, query, id).Scan(
		&discipline.ID,
		&discipline.Name,
		&discipline.TeacherID,
		&discipline.GroupID,
		&discipline.PeriodID,
		&discipline.Description,
		&discipline.CreatedAt,
		&discipline.UpdatedAt,
		&discipline.Teacher.ID,
		&teacherCreatedAt,
		&teacherUpdatedAt,
		&discipline.Teacher.User.Username,
		&discipline.Teacher.User.Name,
		&discipline.Teacher.User.Role,
		&discipline.Group.ID,
		&discipline.Group.Name,
		&groupCreatedAt,
		&discipline.Period.ID,
		&discipline.Period.Name,
		&discipline.Period.HalfYear,
		&discipline.Period.StartDate,
		&discipline.Period.EndDate,
		&periodCreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get discipline: %v", err)
	}

	discipline.Teacher.CreatedAt = teacherCreatedAt
	discipline.Teacher.UpdatedAt = teacherUpdatedAt
	discipline.Teacher.User.ID = discipline.Teacher.ID
	discipline.Group.CreatedAt = groupCreatedAt
	discipline.Period.CreatedAt = periodCreatedAt

	return discipline, nil
}

func (r *DisciplineRepository) GetByTeacherID(ctx context.Context, teacherID int) ([]models.Discipline, error) {
	query := `
SELECT d.id, d.name, d.teacher_id, d.group_id, d.period_id, d.description, 
       d.created_at, d.updated_at,
       g.id, g.name, g.created_at,
       p.id, p.name, p.half_year, p.start_date, p.end_date, p.created_at
FROM data.disciplines d
LEFT JOIN data.groups g ON g.id = d.group_id
LEFT JOIN data.periods p ON p.id = d.period_id
WHERE d.teacher_id = $1
ORDER BY d.name
	`

	rows, err := r.db.Query(ctx, query, teacherID)
	if err != nil {
		return nil, fmt.Errorf("failed to query disciplines: %v", err)
	}
	defer rows.Close()

	var disciplines []models.Discipline

	for rows.Next() {
		var discipline models.Discipline
		var groupCreatedAt time.Time
		var periodCreatedAt time.Time

		err := rows.Scan(
			&discipline.ID,
			&discipline.Name,
			&discipline.TeacherID,
			&discipline.GroupID,
			&discipline.PeriodID,
			&discipline.Description,
			&discipline.CreatedAt,
			&discipline.UpdatedAt,
			&discipline.Group.ID,
			&discipline.Group.Name,
			&groupCreatedAt,
			&discipline.Period.ID,
			&discipline.Period.Name,
			&discipline.Period.HalfYear,
			&discipline.Period.StartDate,
			&discipline.Period.EndDate,
			&periodCreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan discipline: %v", err)
		}

		discipline.Group.CreatedAt = groupCreatedAt
		discipline.Period.CreatedAt = periodCreatedAt

		disciplines = append(disciplines, discipline)
	}

	return disciplines, nil
}

func (r *DisciplineRepository) GetByGroupID(ctx context.Context, groupID int) ([]models.Discipline, error) {
	query := `
SELECT d.id, d.name, d.teacher_id, d.group_id, d.period_id, d.description, 
       d.created_at, d.updated_at,
       t.id, t.created_at, t.updated_at,
       u.username, u.name, u.role,
       p.id, p.name, p.half_year, p.start_date, p.end_date, p.created_at
FROM data.disciplines d
LEFT JOIN data.teachers t ON t.id = d.teacher_id
LEFT JOIN auth.users u ON u.id = t.id
LEFT JOIN data.periods p ON p.id = d.period_id
WHERE d.group_id = $1
ORDER BY d.name
	`

	rows, err := r.db.Query(ctx, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to query disciplines: %v", err)
	}
	defer rows.Close()

	var disciplines []models.Discipline

	for rows.Next() {
		var discipline models.Discipline
		var teacherCreatedAt, teacherUpdatedAt time.Time
		var periodCreatedAt time.Time

		err := rows.Scan(
			&discipline.ID,
			&discipline.Name,
			&discipline.TeacherID,
			&discipline.GroupID,
			&discipline.PeriodID,
			&discipline.Description,
			&discipline.CreatedAt,
			&discipline.UpdatedAt,
			&discipline.Teacher.ID,
			&teacherCreatedAt,
			&teacherUpdatedAt,
			&discipline.Teacher.User.Username,
			&discipline.Teacher.User.Name,
			&discipline.Teacher.User.Role,
			&discipline.Period.ID,
			&discipline.Period.Name,
			&discipline.Period.HalfYear,
			&discipline.Period.StartDate,
			&discipline.Period.EndDate,
			&periodCreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan discipline: %v", err)
		}

		discipline.Teacher.CreatedAt = teacherCreatedAt
		discipline.Teacher.UpdatedAt = teacherUpdatedAt
		discipline.Teacher.User.ID = discipline.Teacher.ID
		discipline.Period.CreatedAt = periodCreatedAt

		disciplines = append(disciplines, discipline)
	}

	return disciplines, nil
}

func (r *DisciplineRepository) Update(ctx context.Context, discipline *models.Discipline) error {
	query := `
UPDATE data.disciplines
SET name = $1, teacher_id = $2, group_id = $3, period_id = $4, 
    description = $5, updated_at = CURRENT_TIMESTAMP
WHERE id = $6
RETURNING updated_at
	`

	return r.db.QueryRow(ctx, query,
		discipline.Name,
		discipline.TeacherID,
		discipline.GroupID,
		discipline.PeriodID,
		discipline.Description,
		discipline.ID,
	).Scan(&discipline.UpdatedAt)
}

func (r *DisciplineRepository) Delete(ctx context.Context, id int) error {
	query := `
DELETE FROM data.disciplines
WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *DisciplineRepository) GetAll(ctx context.Context) ([]models.Discipline, error) {
	query := `
SELECT d.id, d.name, d.teacher_id, d.group_id, d.period_id, d.description, 
       d.created_at, d.updated_at,
       t.id, t.created_at, t.updated_at,
       u.username, u.name, u.role,
       g.id, g.name, g.created_at,
       p.id, p.name, p.half_year, p.start_date, p.end_date, p.created_at
FROM data.disciplines d
LEFT JOIN data.teachers t ON t.id = d.teacher_id
LEFT JOIN auth.users u ON u.id = t.id
LEFT JOIN data.groups g ON g.id = d.group_id
LEFT JOIN data.periods p ON p.id = d.period_id
ORDER BY d.name
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query disciplines: %v", err)
	}
	defer rows.Close()

	var disciplines []models.Discipline

	for rows.Next() {
		var discipline models.Discipline
		var teacherCreatedAt, teacherUpdatedAt time.Time
		var groupCreatedAt time.Time
		var periodCreatedAt time.Time

		err := rows.Scan(
			&discipline.ID,
			&discipline.Name,
			&discipline.TeacherID,
			&discipline.GroupID,
			&discipline.PeriodID,
			&discipline.Description,
			&discipline.CreatedAt,
			&discipline.UpdatedAt,
			&discipline.Teacher.ID,
			&teacherCreatedAt,
			&teacherUpdatedAt,
			&discipline.Teacher.User.Username,
			&discipline.Teacher.User.Name,
			&discipline.Teacher.User.Role,
			&discipline.Group.ID,
			&discipline.Group.Name,
			&groupCreatedAt,
			&discipline.Period.ID,
			&discipline.Period.Name,
			&discipline.Period.HalfYear,
			&discipline.Period.StartDate,
			&discipline.Period.EndDate,
			&periodCreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan discipline: %v", err)
		}

		discipline.Teacher.CreatedAt = teacherCreatedAt
		discipline.Teacher.UpdatedAt = teacherUpdatedAt
		discipline.Teacher.User.ID = discipline.Teacher.ID
		discipline.Group.CreatedAt = groupCreatedAt
		discipline.Period.CreatedAt = periodCreatedAt

		disciplines = append(disciplines, discipline)
	}

	return disciplines, nil
}
