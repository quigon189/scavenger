package services

import (
	"context"
	"fmt"

	"data-service/internal/models"
)

type LabService struct {
	service *Service
}

func NewLabService(service *Service) *LabService {
	return &LabService{service: service}
}

func (l *LabService) CreateLab(ctx context.Context, lab *models.Lab) error {
	// Только преподаватель может создавать лабораторные
	user, err := l.service.RequireRole(ctx, "teacher")
	if err != nil {
		return err
	}

	// Проверяем, что дисциплина принадлежит преподавателю
	discipline, err := l.service.repo.Disciplines.GetByID(ctx, lab.DisciplineID)
	if err != nil {
		return fmt.Errorf("failed to get discipline: %v", err)
	}

	teacher, err := l.service.repo.Teachers.GetByID(ctx, user.ID)
	if err != nil || teacher == nil || discipline.TeacherID != teacher.ID {
		return fmt.Errorf("access denied: discipline does not belong to teacher")
	}

	return l.service.repo.Labs.Create(ctx, lab)
}

func (l *LabService) GetLab(ctx context.Context, id int) (*models.Lab, error) {
	user, err := l.service.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	lab, err := l.service.repo.Labs.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Проверка прав доступа
	switch user.Role {
	case "student":
		student, err := l.service.repo.Students.GetByID(ctx, user.ID)
		if err != nil || student == nil || student.GroupID != lab.Discipline.GroupID {
			return nil, fmt.Errorf("access denied")
		}
	case "teacher":
		teacher, err := l.service.repo.Teachers.GetByID(ctx, user.ID)
		if err != nil || teacher == nil || teacher.ID != lab.Discipline.TeacherID {
			return nil, fmt.Errorf("access denied")
		}
	}

	// Загружаем файлы лабораторной
	files, err := l.service.repo.Labs.GetFiles(ctx, id)
	if err == nil {
		lab.Files = files
	}

	return lab, nil
}

func (l *LabService) GetLabsByDiscipline(ctx context.Context, sessionID string, disciplineID int) ([]models.Lab, error) {
	user, err := l.service.Authenticate(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Проверяем права доступа к дисциплине
	discipline, err := l.service.repo.Disciplines.GetByID(ctx, disciplineID)
	if err != nil {
		return nil, err
	}

	switch user.Role {
	case "student":
		student, err := l.service.repo.Students.GetByUserID(ctx, user.ID)
		if err != nil || student == nil || student.GroupID != discipline.GroupID {
			return nil, fmt.Errorf("access denied")
		}
	case "teacher":
		teacher, err := l.service.repo.Teachers.GetByUserID(ctx, user.ID)
		if err != nil || teacher == nil || teacher.ID != discipline.TeacherID {
			return nil, fmt.Errorf("access denied")
		}
	}

	labs, err := l.service.repo.Labs.GetByDisciplineID(ctx, disciplineID)
	if err != nil {
		return nil, err
	}

	// Загружаем файлы для каждой лабораторной
	for i := range labs {
		files, err := l.service.repo.Labs.GetFiles(ctx, labs[i].ID)
		if err == nil {
			labs[i].Files = files
		}
	}

	return labs, nil
}

func (l *LabService) UpdateLab(ctx context.Context, sessionID string, lab *models.Lab) error {
	user, err := l.service.RequireRole(ctx, sessionID, "teacher")
	if err != nil {
		return err
	}

	// Проверяем, что лабораторная принадлежит преподавателю
	existingLab, err := l.service.repo.Labs.GetByID(ctx, lab.ID)
	if err != nil {
		return fmt.Errorf("failed to get lab: %v", err)
	}

	teacher, err := l.service.repo.Teachers.GetByUserID(ctx, user.ID)
	if err != nil || teacher == nil || existingLab.Discipline.TeacherID != teacher.ID {
		return fmt.Errorf("access denied: lab does not belong to teacher")
	}

	return l.service.repo.Labs.Update(ctx, lab)
}

func (l *LabService) DeleteLab(ctx context.Context, sessionID string, id int) error {
	user, err := l.service.RequireRole(ctx, sessionID, "teacher")
	if err != nil {
		return err
	}

	// Проверяем, что лабораторная принадлежит преподавателю
	lab, err := l.service.repo.Labs.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get lab: %v", err)
	}

	teacher, err := l.service.repo.Teachers.GetByUserID(ctx, user.ID)
	if err != nil || teacher == nil || lab.Discipline.TeacherID != teacher.ID {
		return fmt.Errorf("access denied: lab does not belong to teacher")
	}

	return l.service.repo.Labs.Delete(ctx, id)
}

func (l *LabService) AddFileToLab(ctx context.Context, sessionID string, labID, fileID int) error {
	user, err := l.service.RequireRole(ctx, sessionID, "teacher")
	if err != nil {
		return err
	}

	// Проверяем, что лабораторная принадлежит преподавателю
	lab, err := l.service.repo.Labs.GetByID(ctx, labID)
	if err != nil {
		return fmt.Errorf("failed to get lab: %v", err)
	}

	teacher, err := l.service.repo.Teachers.GetByUserID(ctx, user.ID)
	if err != nil || teacher == nil || lab.Discipline.TeacherID != teacher.ID {
		return fmt.Errorf("access denied: lab does not belong to teacher")
	}

	return l.service.repo.Labs.AddFile(ctx, labID, fileID)
}
