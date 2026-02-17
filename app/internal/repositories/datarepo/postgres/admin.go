package postgres

import (
	"context"
	"scavenger/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepo struct {
	db *pgxpool.Pool
}

func NewAdminRepo(db *pgxpool.Pool) *AdminRepo {
	return &AdminRepo{db: db}
}

func (r *AdminRepo) GetAdminDashboard(ctx context.Context) (*models.AdminDashboard, error) {
	query := `
SELECT 
    (SELECT COUNT(*) FROM data.students),
    (SELECT COUNT(*) FROM data.teachers),
    (SELECT COUNT(*) FROM data.groups),
    (SELECT COUNT(*) FROM data.disciplines),
	(SELECT COUNT(*) FROM data.labs),
	(SELECT COUNT(*) FROM auth.users u WHERE u.status = 'pending' AND u.role = 'student'),
	(SELECT COUNT(*) FROM data.lab_reports r WHERE r.status = 'submitted');
	`

	adminDashboard := &models.AdminDashboard{}
	err := r.db.QueryRow(ctx, query).Scan(
		&adminDashboard.TotalStudents,
		&adminDashboard.TotalTeachers,
		&adminDashboard.TotalGroups,
		&adminDashboard.TotalDisciplines,
		&adminDashboard.TotalLabs,
		&adminDashboard.PendingStudents,
		&adminDashboard.RecentReports,
	)
	if err != nil {
		return nil, err
	}

	return adminDashboard, nil
}
