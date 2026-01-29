package services

import (
	"context"
	"fmt"

	"data-service/internal/models"
)

type ReportService struct {
	service *Service
}

func NewReportService(service *Service) *ReportService {
	return &ReportService{service: service}
}

func (r *ReportService) CreateReport(ctx context.Context, labID int) (*models.LabReport, error) {
	user, err := r.service.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	if user.Role != "student" {
		return nil, fmt.Errorf("only students can create reports")
	}

	// Получаем студента
	student, err := r.service.repo.Students.GetByID(ctx, user.ID)
	if err != nil || student == nil {
		return nil, fmt.Errorf("student not found")
	}

	// Проверяем лабораторную
	lab, err := r.service.repo.Labs.GetByID(ctx, labID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lab: %v", err)
	}

	// Проверяем, что студент в нужной группе
	if student.GroupID != lab.Discipline.GroupID {
		return nil, fmt.Errorf("student is not in the correct group for this lab")
	}

	// Создаем отчет
	report := &models.LabReport{
		LabID:     labID,
		StudentID: student.ID,
		Status:    "submitted",
	}

	err = r.service.repo.Reports.Create(ctx, report)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (r *ReportService) GetReport(ctx context.Context, id int) (*models.LabReport, error) {
	user, err := r.service.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	report, err := r.service.repo.Reports.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Проверка прав доступа
	switch user.Role {
	case "student":
		if report.Student.ID != user.ID {
			return nil, fmt.Errorf("access denied")
		}
	case "teacher":
		teacher, err := r.service.repo.Teachers.GetByID(ctx, user.ID)
		if err != nil || teacher == nil || report.Lab.Discipline.TeacherID != teacher.ID {
			return nil, fmt.Errorf("access denied")
		}
	}

	// Загружаем файлы отчета
	files, err := r.service.repo.Reports.GetFiles(ctx, id)
	if err == nil {
		report.Files = files
	}

	return report, nil
}

func (r *ReportService) GetReportsByLab(ctx context.Context, labID int) ([]models.LabReport, error) {
	user, err := r.service.RequireRole(ctx, "teacher")
	if err != nil {
		return nil, err
	}

	// Проверяем, что лабораторная принадлежит преподавателю
	lab, err := r.service.repo.Labs.GetByID(ctx, labID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lab: %v", err)
	}

	teacher, err := r.service.repo.Teachers.GetByID(ctx, user.ID)
	if err != nil || teacher == nil || lab.Discipline.TeacherID != teacher.ID {
		return nil, fmt.Errorf("access denied: lab does not belong to teacher")
	}

	reports, err := r.service.repo.Reports.GetByLabID(ctx, labID)
	if err != nil {
		return nil, err
	}

	// Загружаем файлы для каждого отчета
	for i := range reports {
		files, err := r.service.repo.Reports.GetFiles(ctx, reports[i].ID)
		if err == nil {
			reports[i].Files = files
		}
	}

	return reports, nil
}

func (r *ReportService) GetReportsByStudent(ctx context.Context, studentID int) ([]models.LabReport, error) {
	user, err := r.service.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	// Студент может видеть только свои отчеты
	if user.Role == "student" {
		student, err := r.service.repo.Students.GetByID(ctx, user.ID)
		if err != nil || student == nil || student.ID != studentID {
			return nil, fmt.Errorf("access denied")
		}
	}

	reports, err := r.service.repo.Reports.GetByStudentID(ctx, studentID)
	if err != nil {
		return nil, err
	}

	// Загружаем файлы для каждого отчета
	for i := range reports {
		files, err := r.service.repo.Reports.GetFiles(ctx, reports[i].ID)
		if err == nil {
			reports[i].Files = files
		}
	}

	return reports, nil
}

func (r *ReportService) GradeReport(ctx context.Context, reportID, grade int, teacherNote string) error {
	user, err := r.service.RequireRole(ctx, "teacher")
	if err != nil {
		return err
	}

	// Проверяем, что отчет принадлежит преподавателю
	report, err := r.service.repo.Reports.GetByID(ctx, reportID)
	if err != nil {
		return fmt.Errorf("failed to get report: %v", err)
	}

	teacher, err := r.service.repo.Teachers.GetByID(ctx, user.ID)
	if err != nil || teacher == nil || report.Lab.Discipline.TeacherID != teacher.ID {
		return fmt.Errorf("access denied: report does not belong to teacher")
	}

	// Обновляем оценку
	err = r.service.repo.Reports.UpdateGrade(ctx, reportID, grade)
	if err != nil {
		return err
	}

	// Обновляем комментарии
	report.TeacherNote = teacherNote
	return r.service.repo.Reports.Update(ctx, report)
}

func (r *ReportService) UpdateComment(ctx context.Context, reportID int, comment string) error {
	user, err := r.service.Authenticate(ctx)
	if err != nil {
		return err
	}

	report, err := r.service.repo.Reports.GetByID(ctx, reportID)
	if err != nil {
		return fmt.Errorf("failed to get report: %v", err)
	}

	if user.ID != report.StudentID {
		return fmt.Errorf("access denied")
	}

	report.Comment = comment

	return r.service.repo.Reports.Update(ctx, report)
}

func (r *ReportService) AddFileToReport(ctx context.Context, reportID, fileID int) error {
	user, err := r.service.Authenticate(ctx)
	if err != nil {
		return err
	}

	report, err := r.service.repo.Reports.GetByID(ctx, reportID)
	if err != nil {
		return fmt.Errorf("failed to get report: %v", err)
	}

	// Только студент, создавший отчет, может добавлять файлы
	if user.Role == "student" && report.Student.ID != user.ID {
		return fmt.Errorf("access denied")
	}

	// Проверяем, что отчет еще не оценен
	if report.Status == "graded" {
		return fmt.Errorf("cannot add files to graded report")
	}

	return r.service.repo.Reports.AddFile(ctx, reportID, fileID)
}
