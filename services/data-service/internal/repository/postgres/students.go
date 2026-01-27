package postgres

import (
	"context"
	"data-service/internal/models"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StudentRepository struct {
	db *pgxpool.Pool
}

func NewStudentRepository(db *pgxpool.Pool) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) Create(ctx context.Context, student *models.Student) error {
	query := `
INSERT INTO data.students (user_id, group_id)
VALUES ($1, $2)
RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(ctx, query,
		student.UserID,
		student.GroupID,
	).Scan(
		&student.ID,
		&student.CreatedAt,
		&student.UpdatedAt,
	)
}

func (r *StudentRepository) GetByID(ctx context.Context, id int) (*models.Student, error) {
	query := `
SELECT s.id, s.user_id, s.group_id, s.created_at, s.updated_at,
       u.username, u.name, u.role,
	   g.name, g.created_at
FROM data.students s
JOIN auth.users u ON u.id = s.user_id
JOIN data.groups g ON g.id = s.group_id
WHERE s.id = $1
	`

	student := &models.Student{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&student.ID,
		&student.UserID,
		&student.CreatedAt,
		&student.UpdatedAt,
		&student.User.Username,
		&student.User.Name,
		&student.User.Role,
		&student.Group.Name,
		&student.Group.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get student: %v", err)
	}

	student.User.ID = student.UserID
	student.Group.ID = student.GroupID

	return student, nil
}

func (r *StudentRepository) GetByUserID(ctx context.Context, userID int) (*models.Student, error) {
	query := `
SELECT s.id, s.user_id, s.group_id, s.created_at, s.updated_at,
       u.username, u.name, u.role,
	   g.name, g.created_at
FROM data.students s
JOIN auth.users u ON u.id = s.user_id
JOIN data.groups g ON g.id = s.group_id
WHERE u.id = $1
	`

	student := &models.Student{}

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&student.ID,
		&student.UserID,
		&student.CreatedAt,
		&student.UpdatedAt,
		&student.User.Username,
		&student.User.Name,
		&student.User.Role,
		&student.Group.Name,
		&student.Group.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get student: %v", err)
	}

	student.User.ID = student.UserID
	student.Group.ID = student.GroupID

	return student, nil
}

func (r *StudentRepository) GetByGroupID(ctx context.Context, groupID int) ([]models.Student, error) {
	query := `
SELECT s.id, s.user_id, s.group_id, s.created_at, s.updated_at,
       u.username, u.name, u.role,
	   g.name, g.created_at
FROM data.students s
JOIN auth.users u ON u.id = s.user_id
JOIN data.groups g ON g.id = s.group_id
WHERE g.id = $1
ORDER BY u.name
	`

	rows, err := r.db.Query(ctx, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to query students: %v", err)
	}
	defer rows.Close()

	var students []models.Student

	for rows.Next() {
		var student models.Student

		err := rows.Scan(
			&student.ID,
			&student.UserID,
			&student.CreatedAt,
			&student.UpdatedAt,
			&student.User.Username,
			&student.User.Name,
			&student.User.Role,
			&student.Group.Name,
			&student.Group.CreatedAt,
		)
		if err != nil {
			if err == pgx.ErrNoRows {
				return nil, nil
			}
			return nil, fmt.Errorf("failed to scan student: %v", err)
		}

		student.User.ID = student.UserID
		student.Group.ID = student.GroupID

		students = append(students, student)
	}

	return students, nil

}

func (r *StudentRepository) Update(ctx context.Context, student *models.Student) error {
	query := `
UPDATE data.students
SET group_id = $1, updated_at = CURRENT_TIMESTAMP
WHERE id = $2
RETURNING updated_at
	`

	return r.db.QueryRow(ctx, query, student.GroupID, student.ID).Scan(&student.UpdatedAt)
}

func (r *StudentRepository) Delete(ctx context.Context, id int) error {
	query := `
DELETE FROM data.students
WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *StudentRepository) GetAll(ctx context.Context) ([]models.Student, error) {
	query := `
SELECT s.id, s.user_id, s.group_id, s.created_at, s.updated_at,
       u.username, u.name, u.role,
	   g.name, g.created_at
FROM data.students s
JOIN auth.users u ON u.id = s.user_id
JOIN data.groups g ON g.id = s.group_id
ORDER BY u.name
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query students: %v", err)
	}
	defer rows.Close()

	var students []models.Student

	for rows.Next() {
		var student models.Student

		err := rows.Scan(
			&student.ID,
			&student.UserID,
			&student.CreatedAt,
			&student.UpdatedAt,
			&student.User.Username,
			&student.User.Name,
			&student.User.Role,
			&student.Group.Name,
			&student.Group.CreatedAt,
		)
		if err != nil {
			if err == pgx.ErrNoRows {
				return nil, nil
			}
			return nil, fmt.Errorf("failed to scan student: %v", err)
		}

		student.User.ID = student.UserID
		student.Group.ID = student.GroupID

		students = append(students, student)
	}

	return students, nil
}
