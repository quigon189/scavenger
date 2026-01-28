package services

import (
	"context"
	"fmt"

	"data-service/internal/models"
)

type StudentService struct {
	service *Service
}

func NewStudentService(service *Service) *StudentService {
	return &StudentService{service: service}
}

func (s *StudentService) CreateStudent(ctx context.Context, student *models.Student) error {
	return s.service.repo.Students.Create(ctx, student)
}

func (s *StudentService) GetStudent(ctx context.Context, id int) (*models.Student, error) {
	user, err := s.service.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	student, err := s.service.repo.Students.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Студент может видеть только себя
	if user.Role == "student" && student.ID != user.ID {
		return nil, fmt.Errorf("access denied")
	}

	return student, nil
}

func (s *StudentService) UpdateStudent(ctx context.Context, student *models.Student) error {
	// Только админ может обновлять студентов
	_, err := s.service.RequireRole(ctx, "admin")
	if err != nil {
		return err
	}

	return s.service.repo.Students.Update(ctx, student)
}

func (s *StudentService) DeleteStudent(ctx context.Context, id int) error {
	// Только админ может удалять студентов
	_, err := s.service.RequireRole(ctx, "admin")
	if err != nil {
		return err
	}

	return s.service.repo.Students.Delete(ctx, id)
}

func (s *StudentService) GetStudentsByGroup(ctx context.Context, groupID int) ([]models.Student, error) {
	// Только преподаватели и админы могут видеть список студентов
	_, err := s.service.RequireRole(ctx, "admin", "teacher")
	if err != nil {
		return nil, err
	}

	return s.service.repo.Students.GetByGroupID(ctx, groupID)
}

func (s *StudentService) GetAllStudents(ctx context.Context) ([]models.Student, error) {
	// Только админ может видеть всех студентов
	_, err := s.service.RequireRole(ctx, "admin")
	if err != nil {
		return nil, err
	}

	return s.service.repo.Students.GetAll(ctx)
}
