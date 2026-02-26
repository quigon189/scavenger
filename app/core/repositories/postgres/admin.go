package postgres

import (
	"context"
	"scavenger/core/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type adminRepo struct {
	db *pgxpool.Pool
}

func NewAdminRepo(db *pgxpool.Pool) *adminRepo {
	return &adminRepo{db: db}
}

func (r *adminRepo) GetAdminDashboard(ctx context.Context) (*models.AdminDashboard, error) {
	query := `
SELECT 
    (SELECT COUNT(*) FROM data.students),
    (SELECT COUNT(*) FROM data.teachers),
    (SELECT COUNT(*) FROM data.groups),
    (SELECT COUNT(*) FROM data.disciplines),
	(SELECT COUNT(*) FROM data.labs),
	(SELECT COUNT(*) FROM auth.registration_codes)
	`

	adminDashboard := &models.AdminDashboard{}
	err := r.db.QueryRow(ctx, query).Scan(
		&adminDashboard.TotalStudents,
		&adminDashboard.TotalTeachers,
		&adminDashboard.TotalGroups,
		&adminDashboard.TotalDisciplines,
		&adminDashboard.TotalLabs,
		&adminDashboard.TotalCodes,
	)
	if err != nil {
		return nil, err
	}

	return adminDashboard, nil
}
