package postgres

import (
	"scavenger/core/repositories"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDataRepo(db *pgxpool.Pool) *repositories.DataRepo {
	return &repositories.DataRepo{
		Admin:      NewAdminRepo(db),
		Student:    NewStudentRepository(db),
		Teacher:    NewTeacherRepository(db),
		Period:     NewPeriodRepository(db),
		Discipline: NewDisciplineRepository(db),
		Lab:        NewLabRepository(db),
		Report:     NewReportRepository(db),
		File:       NewFileRepository(db),
		Group:      NewGroupRepository(db),
	}
}
