package dataservice

import (
	"context"
	"fmt"

	"scavenger/internal/models"
)

func (s *DataService) CreateTeacher(ctx context.Context, teacher *models.Teacher) error {
	// Только админ может создавать преподавателей
	_, err := s.RequireRole(ctx, "admin")
	if err != nil {
		return err
	}

	return s.dataRepo.Teachers.Create(ctx, teacher)
}

func (s *DataService) GetTeacher(ctx context.Context, id int) (*models.Teacher, error) {
	_, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	return s.dataRepo.Teachers.GetByID(ctx, id)
}

func (s *DataService) GetTeacherWithDisciplines(ctx context.Context, id int) (*models.Teacher, error) {
	teacher, err := s.GetTeacher(ctx, id)
	if err != nil {
		return nil, err
	}

	disciplines, err := s.dataRepo.Disciplines.GetByTeacherID(ctx, id)
	if err != nil {
		return teacher, fmt.Errorf("failed to get teacher disciplines: %v", err)
	}

	teacher.Disciplines = disciplines
	return teacher, nil
}

func (s *DataService) GetAllTeachers(ctx context.Context) ([]models.Teacher, error) {
	_, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	return s.dataRepo.Teachers.GetAll(ctx)
}

func (s *DataService) DeleteTeacher(ctx context.Context, id int) error {
	_, err := s.RequireRole(ctx, "admin")
	if err != nil {
		return err
	}

	return s.dataRepo.Teachers.Delete(ctx, id)
}
