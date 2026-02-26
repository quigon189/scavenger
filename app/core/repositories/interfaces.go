package repositories

import (
	"context"

	"scavenger/core/models"
)

type AdminRepo interface {
	GetAdminDashboard(ctx context.Context) (*models.AdminDashboard, error)
}

type AuthRepo interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByID(ctx context.Context, userID int) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	UpdatePassword(ctx context.Context, userID int, passwordHash string) error
	DeleteUser(ctx context.Context, userID int) error

	CreateCode(ctx context.Context, code *models.RegistrationCode) error
	GetCode(ctx context.Context, codeStr string) (*models.RegistrationCode, error)
	GetAllCodes(ctx context.Context) ([]models.RegistrationCode, error)
	DeleteCode(ctx context.Context, code string) error
}

type DisciplineRepo interface {
	Create(ctx context.Context, discipline *models.Discipline) error
	GetByID(ctx context.Context, id int) (*models.Discipline, error)
	GetByTeacherID(ctx context.Context, teacherID int) ([]models.Discipline, error)
	GetByGroupID(ctx context.Context, groupID int) ([]models.Discipline, error)
	Update(ctx context.Context, discipline *models.Discipline) error
	Delete(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]models.Discipline, error)
}

type FileRepo interface {
	Create(ctx context.Context, file *models.File) error
	GetByID(ctx context.Context, id int) (*models.File, error)
	GetByUUID(ctx context.Context, uuid string) (*models.File, error)
	Delete(ctx context.Context, id int) error
	DeleteByUUID(ctx context.Context, uuid string) error
}

type GroupRepo interface {
	Create(ctx context.Context, group *models.Group) error
	GetByID(ctx context.Context, id int) (*models.Group, error)
	GetByName(ctx context.Context, name string) (*models.Group, error)
	GetWithStudents(ctx context.Context, id int) (*models.Group, error)
	Update(ctx context.Context, group *models.Group) error
	Delete(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]models.Group, error)
}

type LabRepo interface {
	Create(ctx context.Context, lab *models.Lab) error
	GetByID(ctx context.Context, id int) (*models.Lab, error)
	GetByDisciplineID(ctx context.Context, disciplineID int) ([]models.Lab, error)
	Update(ctx context.Context, lab *models.Lab) error
	Delete(ctx context.Context, id int) error
	GetFiles(ctx context.Context, labID int) ([]models.File, error)
	AddFile(ctx context.Context, labID, fileID int) error
	RemoveFile(ctx context.Context, labID, fileID int) error
}

type PeriodRepo interface {
	Create(ctx context.Context, period *models.Period) error
	GetByID(ctx context.Context, id int) (*models.Period, error)
	GetActive(ctx context.Context) (*models.Period, error)
	Update(ctx context.Context, period *models.Period) error
	Delete(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]models.Period, error)
}

type ReportRepo interface {
	Create(ctx context.Context, report *models.LabReport) error
	GetByID(ctx context.Context, id int) (*models.LabReport, error)
	GetByLabID(ctx context.Context, labID int) ([]models.LabReport, error)
	GetByStudentID(ctx context.Context, studentID int) ([]models.LabReport, error)
	Update(ctx context.Context, report *models.LabReport) error
	UpdateStatus(ctx context.Context, id int, status string) error
	UpdateGrade(ctx context.Context, id, grade int) error
	Delete(ctx context.Context, id int) error
	GetFiles(ctx context.Context, reportID int) ([]models.File, error)
	AddFile(ctx context.Context, reportID, fileID int) error
	RemoveFile(ctx context.Context, reportID, fileID int) error
}

type SessionRepo interface {
	SetSession(ctx context.Context, sessionID string, session *models.UserSession) error
	GetSession(ctx context.Context, sessionID string) (*models.UserSession, error)
	DeleteSession(ctx context.Context, sessionID string) error
	RefreshSession(ctx context.Context, sessionID string, session *models.UserSession) error
	GetUserSessions(ctx context.Context, userID int) ([]string, error)
}

type StudentRepo interface {
	Create(ctx context.Context, student *models.Student) error
	GetByID(ctx context.Context, id int) (*models.Student, error)
	GetByGroupID(ctx context.Context, groupID int) ([]models.Student, error)
	Update(ctx context.Context, student *models.Student) error
	Delete(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]models.Student, error)
}

type TeacherRepo interface {
	Create(ctx context.Context, teacher *models.Teacher) error
	GetByID(ctx context.Context, id int) (*models.Teacher, error)
	GetAll(ctx context.Context) ([]models.Teacher, error)
	Delete(ctx context.Context, id int) error
}

type DataRepo struct {
	Student    StudentRepo
	Teacher    TeacherRepo
	Discipline DisciplineRepo
	Lab        LabRepo
	Report     ReportRepo
	Group      GroupRepo
	File       FileRepo
	Period     PeriodRepo
	Admin      AdminRepo
}
