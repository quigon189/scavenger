package services

import (
	"context"
	"fmt"

	"data-service/internal/models"
)

type TeacherService struct {
	service *Service
}

func NewTeacherService(service *Service) *TeacherService {
	return &TeacherService{service: service}
}

func (t *TeacherService) CreateTeacher(ctx context.Context, teacher *models.Teacher) error {
	// Только админ может создавать преподавателей
	_, err := t.service.RequireRole(ctx, "admin")
	if err != nil {
		return err
	}

	return t.service.repo.Teachers.Create(ctx, teacher)
}

func (t *TeacherService) GetTeacher(ctx context.Context, id int) (*models.Teacher, error) {
	_, err := t.service.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	return t.service.repo.Teachers.GetByID(ctx, id)
}

func (t *TeacherService) GetTeacherWithDisciplines(ctx context.Context, id int) (*models.Teacher, error) {
	teacher, err := t.GetTeacher(ctx, id)
	if err != nil {
		return nil, err
	}

	disciplines, err := t.service.repo.Disciplines.GetByTeacherID(ctx, id)
	if err != nil {
		return teacher, fmt.Errorf("failed to get teacher disciplines: %v", err)
	}

	teacher.Disciplines = disciplines
	return teacher, nil
}

func (t *TeacherService) GetAllTeachers(ctx context.Context, sessionID string) ([]models.Teacher, error) {
	_, err := t.service.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	return t.service.repo.Teachers.GetAll(ctx)
}
