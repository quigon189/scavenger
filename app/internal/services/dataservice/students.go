package dataservice

import (
	"context"
	"fmt"

	"scavenger/internal/models"
)

func (s *DataService) CreateStudent(ctx context.Context, student *models.Student) error {
	return s.dataRepo.Students.Create(ctx, student)
}

func (s *DataService) GetStudent(ctx context.Context, id int) (*models.Student, error) {
	user, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	student, err := s.dataRepo.Students.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Студент может видеть только себя
	if user.Role == "student" && student.ID != user.ID {
		return nil, fmt.Errorf("access denied")
	}

	return student, nil
}

func (s *DataService) UpdateStudent(ctx context.Context, student *models.Student) error {
	// Только админ может обновлять студентов
	_, err := s.RequireRole(ctx, "admin")
	if err != nil {
		return err
	}

	return s.dataRepo.Students.Update(ctx, student)
}

func (s *DataService) DeleteStudent(ctx context.Context, id int) error {
	// Только админ может удалять студентов
	_, err := s.RequireRole(ctx, "admin")
	if err != nil {
		return err
	}

	return s.dataRepo.Students.Delete(ctx, id)
}

func (s *DataService) GetStudentsByGroup(ctx context.Context, groupID int) ([]models.Student, error) {
	// Только преподаватели и админы могут видеть список студентов
	_, err := s.RequireRole(ctx, "admin", "teacher")
	if err != nil {
		return nil, err
	}

	return s.dataRepo.Students.GetByGroupID(ctx, groupID)
}

func (s *DataService) GetAllStudents(ctx context.Context) ([]models.Student, error) {
	// Только админ может видеть всех студентов
	_, err := s.RequireRole(ctx, "admin")
	if err != nil {
		return nil, err
	}

	return s.dataRepo.Students.GetAll(ctx)
}
