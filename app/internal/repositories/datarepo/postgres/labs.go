package postgres

import (
	"context"
	"fmt"

	"scavenger/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LabRepository struct {
	db *pgxpool.Pool
}

func NewLabRepository(db *pgxpool.Pool) *LabRepository {
	return &LabRepository{db: db}
}

func (r *LabRepository) Create(ctx context.Context, lab *models.Lab) error {
	query := `
INSERT INTO data.labs (discipline_id, md_content, name, description, deadline)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(ctx, query,
		lab.DisciplineID,
		lab.MDContent,
		lab.Name,
		lab.Description,
		lab.Deadline,
	).Scan(
		&lab.ID,
		&lab.CreatedAt,
		&lab.UpdatedAt,
	)
}

func (r *LabRepository) GetByID(ctx context.Context, id int) (*models.Lab, error) {
	query := `
SELECT l.id, l.discipline_id, l.md_content, l.name, l.description, l.deadline,
       l.created_at, l.updated_at,
       d.id, d.name, d.teacher_id, d.group_id, d.period_id, d.description,
       d.created_at, d.updated_at 
FROM data.labs l
JOIN data.disciplines d ON d.id = l.discipline_id
WHERE l.id = $1
	`

	lab := &models.Lab{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&lab.ID,
		&lab.DisciplineID,
		&lab.MDContent,
		&lab.Name,
		&lab.Description,
		&lab.Deadline,
		&lab.CreatedAt,
		&lab.UpdatedAt,
		&lab.Discipline.ID,
		&lab.Discipline.Name,
		&lab.Discipline.TeacherID,
		&lab.Discipline.GroupID,
		&lab.Discipline.PeriodID,
		&lab.Discipline.Description,
		&lab.Discipline.CreatedAt,
		&lab.Discipline.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get lab: %v", err)
	}

	return lab, nil
}

func (r *LabRepository) GetByDisciplineID(ctx context.Context, disciplineID int) ([]models.Lab, error) {
	query := `
SELECT l.id, l.discipline_id, l.md_content, l.name, l.description, l.deadline,
       l.created_at, l.updated_at
FROM data.labs l
WHERE l.discipline_id = $1
ORDER BY l.deadline
	`

	rows, err := r.db.Query(ctx, query, disciplineID)
	if err != nil {
		return nil, fmt.Errorf("failed to query labs: %v", err)
	}
	defer rows.Close()

	var labs []models.Lab

	for rows.Next() {
		var lab models.Lab

		err := rows.Scan(
			&lab.ID,
			&lab.DisciplineID,
			&lab.MDContent,
			&lab.Name,
			&lab.Description,
			&lab.Deadline,
			&lab.CreatedAt,
			&lab.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan lab: %v", err)
		}

		labs = append(labs, lab)
	}

	return labs, nil
}

func (r *LabRepository) Update(ctx context.Context, lab *models.Lab) error {
	query := `
UPDATE data.labs
SET discipline_id = $1, md_content = $2, name = $3, 
    description = $4, deadline = $5, updated_at = CURRENT_TIMESTAMP
WHERE id = $6
RETURNING updated_at
	`

	return r.db.QueryRow(ctx, query,
		lab.DisciplineID,
		lab.MDContent,
		lab.Name,
		lab.Description,
		lab.Deadline,
		lab.ID,
	).Scan(&lab.UpdatedAt)
}

func (r *LabRepository) Delete(ctx context.Context, id int) error {
	query := `
DELETE FROM data.labs
WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *LabRepository) GetFiles(ctx context.Context, labID int) ([]models.File, error) {
	query := `
SELECT f.id, f.uuid, f.filename, f.size, f.content_type, f.bucket, f.path, f.created_at
FROM data.lab_files lf
JOIN data.files f ON f.id = lf.file_id
WHERE lf.lab_id = $1
ORDER BY f.created_at
	`

	rows, err := r.db.Query(ctx, query, labID)
	if err != nil {
		return nil, fmt.Errorf("failed to query lab files: %v", err)
	}
	defer rows.Close()

	var files []models.File

	for rows.Next() {
		var file models.File

		err := rows.Scan(
			&file.ID,
			&file.UUID,
			&file.Filename,
			&file.Size,
			&file.ContentType,
			&file.Bucket,
			&file.Path,
			&file.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan file: %v", err)
		}

		files = append(files, file)
	}

	return files, nil
}

func (r *LabRepository) AddFile(ctx context.Context, labID, fileID int) error {
	query := `
INSERT INTO data.lab_files (lab_id, file_id)
VALUES ($1, $2)
ON CONFLICT (lab_id, file_id) DO NOTHING
	`

	_, err := r.db.Exec(ctx, query, labID, fileID)
	return err
}

func (r *LabRepository) RemoveFile(ctx context.Context, labID, fileID int) error {
	query := `
DELETE FROM data.lab_files
WHERE lab_id = $1 AND file_id = $2
	`

	_, err := r.db.Exec(ctx, query, labID, fileID)
	return err
}
