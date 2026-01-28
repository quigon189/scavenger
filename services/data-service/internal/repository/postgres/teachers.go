package postgres

import (
	"context"
	"data-service/internal/models"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeacherRepository struct {
	db *pgxpool.Pool
}

func NewTeacherRepository(db *pgxpool.Pool) *TeacherRepository {
	return &TeacherRepository{db: db}
}

func (r *TeacherRepository) Create(ctx context.Context, teacher *models.Teacher) error {
	query := `
INSERT INTO data.teachers (id)
VALUES ($1)
RETURNING created_at, updated_at
	`

	return r.db.QueryRow(ctx, query, teacher.ID).Scan(
		&teacher.CreatedAt,
		&teacher.UpdatedAt,
	)
}

func (r *TeacherRepository) GetByID(ctx context.Context, id int) (*models.Teacher, error) {
	query := `
SELECT t.id, t.created_at, t.updated_at,
       u.username, u.name, u.role
FROM data.teachers t
JOIN auth.users u ON u.id = t.id
WHERE t.id = $1
	`

	teacher := &models.Teacher{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&teacher.ID,
		&teacher.CreatedAt,
		&teacher.UpdatedAt,
		&teacher.User.Username,
		&teacher.User.Name,
		&teacher.User.Role,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get teacher: %v", err)
	}

	teacher.User.ID = teacher.ID

	return teacher, nil
}

func (r *TeacherRepository) GetAll(ctx context.Context) ([]models.Teacher, error) {
	query := `
SELECT t.id, t.created_at, t.updated_at,
       u.username, u.name, u.role
FROM data.teachers t
JOIN auth.users u ON u.id = t.id
ORDER BY u.name
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query teachers: %v", err)
	}
	defer rows.Close()

	var teachers []models.Teacher

	for rows.Next() {
		var teacher models.Teacher

		err := rows.Scan(
			&teacher.ID,
			&teacher.CreatedAt,
			&teacher.UpdatedAt,
			&teacher.User.Username,
			&teacher.User.Name,
			&teacher.User.Role,
		)
		if err != nil {
			if err == pgx.ErrNoRows {
				return nil, nil
			}
			return nil, fmt.Errorf("failed to scan student: %v", err)
		}

		teacher.User.ID = teacher.ID

		teachers = append(teachers, teacher)
	}

	return teachers, nil
}

func (r *TeacherRepository) GetWithDisciplines(ctx context.Context, id int) (*models.Teacher, error) {
	teacher, err := r.GetByID(ctx, id)
	if err != nil || teacher == nil {
		return teacher, err
	}

	disciplineRepo := NewDisciplineRepository(r.db)
	disciplines, err := disciplineRepo.GetByTeacherID(ctx, id)
	if err != nil {
		return teacher, fmt.Errorf("failed to get teacher disciplines: %v", err)
	}

	teacher.Disciplines = disciplines
	return teacher, nil
}

func (r *TeacherRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM data.teachers WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	return err
}
