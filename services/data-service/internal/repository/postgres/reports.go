package postgres

import (
	"context"
	"fmt"

	"data-service/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReportRepository struct {
	db *pgxpool.Pool
}

func NewReportRepository(db *pgxpool.Pool) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) Create(ctx context.Context, report *models.LabReport) error {
	query := `
INSERT INTO data.lab_reports (lab_id, student_id)
VALUES ($1, $2)
RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(ctx, query,
		report.LabID,
		report.StudentID,
	).Scan(
		&report.ID,
		&report.CreatedAt,
		&report.UpdatedAt,
	)
}

func (r *ReportRepository) GetByID(ctx context.Context, id int) (*models.LabReport, error) {
	query := `
SELECT lr.id, lr.lab_id, lr.student_id, lr.status, lr.grade, 
       lr.comment, lr.teacher_note, lr.graded_at, lr.created_at, lr.updated_at,
       l.id, l.discipline_id, l.md_file_id, l.name, l.description, l.deadline,
       l.created_at, l.updated_at,
       s.id, s.user_id, s.group_id, s.created_at, s.updated_at,
       u.username, u.name, u.role
FROM data.lab_reports lr
JOIN data.labs l ON l.id = lr.lab_id
JOIN data.students s ON s.id = lr.student_id
JOIN auth.users u ON u.id = s.user_id
WHERE lr.id = $1
	`

	report := &models.LabReport{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&report.ID,
		&report.LabID,
		&report.StudentID,
		&report.Status,
		&report.Grade,
		&report.Comment,
		&report.TeacherNote,
		&report.GradedAt,
		&report.CreatedAt,
		&report.UpdatedAt,
		&report.Lab.ID,
		&report.Lab.DisciplineID,
		&report.Lab.MDFileID,
		&report.Lab.Name,
		&report.Lab.Description,
		&report.Lab.Deadline,
		&report.Lab.CreatedAt,
		&report.Lab.UpdatedAt,
		&report.Student.ID,
		&report.Student.UserID,
		&report.Student.GroupID,
		&report.Student.CreatedAt,
		&report.Student.UpdatedAt,
		&report.Student.User.Username,
		&report.Student.User.Name,
		&report.Student.User.Role,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get report: %v", err)
	}

	report.Student.User.ID = report.Student.UserID

	return report, nil
}

func (r *ReportRepository) GetByLabID(ctx context.Context, labID int) ([]models.LabReport, error) {
	query := `
SELECT lr.id, lr.lab_id, lr.student_id, lr.status, lr.grade, 
       lr.comment, lr.teacher_note, lr.graded_at,
	   lr.created_at, lr.updated_at,
       s.id, s.user_id, s.group_id, s.created_at, s.updated_at,
       u.username, u.name, u.role
FROM data.lab_reports lr
JOIN data.students s ON s.id = lr.student_id
JOIN auth.users u ON u.id = s.user_id
WHERE lr.lab_id = $1
ORDER BY lr.created_at
	`

	rows, err := r.db.Query(ctx, query, labID)
	if err != nil {
		return nil, fmt.Errorf("failed to query reports: %v", err)
	}
	defer rows.Close()

	var reports []models.LabReport

	for rows.Next() {
		var report models.LabReport

		err := rows.Scan(
			&report.ID,
			&report.LabID,
			&report.StudentID,
			&report.Status,
			&report.Grade,
			&report.Comment,
			&report.TeacherNote,
			&report.GradedAt,
			&report.CreatedAt,
			&report.UpdatedAt,
			&report.Student.ID,
			&report.Student.UserID,
			&report.Student.GroupID,
			&report.Student.CreatedAt,
			&report.Student.UpdatedAt,
			&report.Student.User.Username,
			&report.Student.User.Name,
			&report.Student.User.Role,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan report: %v", err)
		}

		report.Student.User.ID = report.Student.UserID

		reports = append(reports, report)
	}

	return reports, nil
}

func (r *ReportRepository) GetByStudentID(ctx context.Context, studentID int) ([]models.LabReport, error) {
	query := `
SELECT lr.id, lr.lab_id, lr.student_id, lr.status, lr.grade, 
       lr.comment, lr.teacher_note, lr.graded_at,
	   lr.created_at, lr.updated_at,
       l.id, l.discipline_id, l.md_file_id, l.name, l.description, l.deadline,
       l.created_at, l.updated_at,
       d.name as discipline_name
FROM data.lab_reports lr
JOIN data.labs l ON l.id = lr.lab_id
JOIN data.disciplines d ON d.id = l.discipline_id
WHERE lr.student_id = $1
ORDER BY lr.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query reports: %v", err)
	}
	defer rows.Close()

	var reports []models.LabReport

	for rows.Next() {
		var report models.LabReport

		err := rows.Scan(
			&report.ID,
			&report.LabID,
			&report.StudentID,
			&report.Status,
			&report.Grade,
			&report.Comment,
			&report.TeacherNote,
			&report.GradedAt,
			&report.CreatedAt,
			&report.UpdatedAt,
			&report.Lab.ID,
			&report.Lab.DisciplineID,
			&report.Lab.MDFileID,
			&report.Lab.Name,
			&report.Lab.Description,
			&report.Lab.Deadline,
			&report.Lab.CreatedAt,
			&report.Lab.UpdatedAt,
			&report.Lab.Discipline.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan report: %v", err)
		}

		reports = append(reports, report)
	}

	return reports, nil
}

func (r *ReportRepository) Update(ctx context.Context, report *models.LabReport) error {
	query := `
UPDATE data.lab_reports
SET comment = $1, teacher_note = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $3
RETURNING updated_at
	`

	return r.db.QueryRow(ctx, query,
		report.Comment,
		report.TeacherNote,
		report.ID,
	).Scan(&report.UpdatedAt)
}

func (r *ReportRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	query := `
UPDATE data.lab_reports
SET status = $1, updated_at = CURRENT_TIMESTAMP
WHERE id = $2
	`

	_, err := r.db.Exec(ctx, query, status, id)
	return err
}

func (r *ReportRepository) UpdateGrade(ctx context.Context, id, grade int) error {
	query := `
UPDATE data.lab_reports
SET grade = $1, graded_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $2
	`

	_, err := r.db.Exec(ctx, query, grade, id)
	if err != nil {
		return err
	}

	return r.UpdateStatus(ctx, id, "graded")

}

func (r *ReportRepository) Delete(ctx context.Context, id int) error {
	query := `
DELETE FROM data.lab_reports
WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *ReportRepository) GetFiles(ctx context.Context, reportID int) ([]models.File, error) {
	query := `
SELECT f.id, f.uuid, f.filename, f.size, f.content_type, f.bucket, f.path, f.created_at
FROM data.report_files rf
JOIN data.files f ON f.id = rf.file_id
WHERE rf.report_id = $1
ORDER BY f.created_at
	`

	rows, err := r.db.Query(ctx, query, reportID)
	if err != nil {
		return nil, fmt.Errorf("failed to query report files: %v", err)
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

func (r *ReportRepository) AddFile(ctx context.Context, reportID, fileID int) error {
	query := `
INSERT INTO data.report_files (report_id, file_id)
VALUES ($1, $2)
ON CONFLICT (report_id, file_id) DO NOTHING
	`

	_, err := r.db.Exec(ctx, query, reportID, fileID)
	return err
}

func (r *ReportRepository) RemoveFile(ctx context.Context, reportID, fileID int) error {
	query := `
DELETE FROM data.report_files
WHERE report_id = $1 AND file_id = $2
	`

	_, err := r.db.Exec(ctx, query, reportID, fileID)
	return err
}
