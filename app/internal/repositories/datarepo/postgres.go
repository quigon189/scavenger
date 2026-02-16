package datarepo

import (
	"scavenger/internal/repositories"
	"scavenger/internal/repositories/datarepo/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDataRepository(db *pgxpool.Pool) *repositories.DataRepository {
	return &repositories.DataRepository{
		Students:    postgres.NewStudentRepository(db),
		Teachers:    postgres.NewTeacherRepository(db),
		Disciplines: postgres.NewDisciplineRepository(db),
		Labs:        postgres.NewLabRepository(db),
		Reports:     postgres.NewReportRepository(db),
		Files:       postgres.NewFileRepository(db),
		Groups:      postgres.NewGroupRepository(db),
		Periods:     postgres.NewPeriodRepository(db),
	}
}
