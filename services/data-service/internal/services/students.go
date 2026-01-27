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

func (s *StudentService) GetStudent(ctx context.Context, sessionID string, id int) (*models.Student, error) {
	user, err := s.service.Authenticate(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Проверяем права доступа
	student, err := s.service.repo.Students.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Студент может видеть только себя
	if user.Role == "student" && student.UserID != user.ID {
		return nil, fmt.Errorf("access denied")
	}

	// Преподаватель может видеть студентов своих дисциплин
	if user.Role == "teacher" {
		teacher, err := s.service.repo.Teachers.GetByUserID(ctx, user.ID)
		if err != nil || teacher == nil {
			return nil, fmt.Errorf("access denied")
		}

		disciplines, err := s.service.repo.Disciplines.GetByTeacherID(ctx, teacher.ID)
		if err != nil {
			return nil, fmt.Errorf("access denied")
		}

		// TODO:
	}

	return student, nil
}

func (s *StudentService) GetStudentByUserID(ctx context.Context, sessionID string, userID int) (*models.Student, error) {
	user, err := s.service.Authenticate(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Пользователь может получить только свою информацию
	if user.Role == "student" && user.ID != userID {
		return nil, fmt.Errorf("access denied")
	}

	return s.service.repo.Students.GetByUserID(ctx, userID)
}

func (s *StudentService) UpdateStudent(ctx context.Context, sessionID string, student *models.Student) error {
	// Только админ может обновлять студентов
	_, err := s.service.RequireRole(ctx, sessionID, "admin")
	if err != nil {
		return err
	}

	return s.service.repo.Students.Update(ctx, student)
}

func (s *StudentService) DeleteStudent(ctx context.Context, sessionID string, id int) error {
	// Только админ может удалять студентов
	_, err := s.service.RequireRole(ctx, sessionID, "admin")
	if err != nil {
		return err
	}

	return s.service.repo.Students.Delete(ctx, id)
}

func (s *StudentService) GetStudentsByGroup(ctx context.Context, sessionID string, groupID int) ([]models.Student, error) {
	// Только преподаватели и админы могут видеть список студентов
	user, err := s.service.Authenticate(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if user.Role != "teacher" && user.Role != "admin" {
		return nil, fmt.Errorf("access denied")
	}

	return s.service.repo.Students.GetByGroupID(ctx, groupID)
}

func (s *StudentService) GetAllStudents(ctx context.Context, sessionID string) ([]models.Student, error) {
	// Только админ может видеть всех студентов
	_, err := s.service.RequireRole(ctx, sessionID, "admin")
	if err != nil {
		return nil, err
	}

	return s.service.repo.Students.GetAll(ctx)
}
