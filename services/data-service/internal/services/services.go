package services

import (
	"data-service/internal/repository/postgres"
	"data-service/internal/session"
	"data-service/internal/storage"
)

type Services struct {
	Main        *Service
	Students    *StudentService
	Teachers    *TeacherService
	Disciplines *DisciplineService
	Labs        *LabService
	Reports     *ReportService
	Files       *FileService
	Groups      *GroupService
	Periods     *PeriodService
}

func NewServices(
	repo *postgres.Repository,
	sessionService *session.SessionService,
	storage storage.Storage,
	bucketName string,
) *Services {
	mainService := NewService(repo, sessionService, storage, bucketName)

	return &Services{
		Main:        mainService,
		Students:    NewStudentService(mainService),
		Teachers:    NewTeacherService(mainService),
		Disciplines: NewDisciplineService(mainService),
		Labs:        NewLabService(mainService),
		Reports:     NewReportService(mainService),
		Files:       NewFileService(mainService, storage),
		Groups:      NewGroupService(mainService),
		Periods:     NewPeriodService(mainService),
	}
}
