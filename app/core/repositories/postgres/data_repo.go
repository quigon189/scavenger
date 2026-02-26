package postgres

import (
	"scavenger/core/repositories/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDataRepo(db *pgxpool.Pool) *DataRepo {
	return &DataRepo{
		Admin: postgres.NewAdminRepo(db),
		Student: postgres.NewStudentRepository(db),
		Teacher: postgres.NewTeacherRepository(db),
		Period: postgres.NewPeriodRepository(db),
		Discipline: postgres.NewDisciplineRepository(db),
		Lab: postgres.NewLabRepository(db),
		Report: postgres.NewReportRepository(db),
		File: postgres.NewFileRepository(db),
		Group: postgres.NewGroupRepository(db),
	}
}
