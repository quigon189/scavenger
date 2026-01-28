package services

import (
	"context"
	"fmt"

	"data-service/internal/models"
)

type DisciplineService struct {
	service *Service
}

func NewDisciplineService(service *Service) *DisciplineService {
	return &DisciplineService{service: service}
}

func (d *DisciplineService) CreateDiscipline(ctx context.Context, discipline *models.Discipline) error {
	_, err := d.service.RequireRole(ctx, "teacher", "admin")
	if err != nil {
		return err
	}

	return d.service.repo.Disciplines.Create(ctx, discipline)
}

func (d *DisciplineService) GetDiscipline(ctx context.Context, id int) (*models.Discipline, error) {
	user, err := d.service.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	discipline, err := d.service.repo.Disciplines.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Проверка прав доступа
	switch user.Role {
	case "student":
		student, err := d.service.repo.Students.GetByID(ctx, user.ID)
		if err != nil || student == nil || student.GroupID != discipline.GroupID {
			return nil, fmt.Errorf("access denied")
		}
	case "teacher":
		teacher, err := d.service.repo.Teachers.GetByID(ctx, user.ID)
		if err != nil || teacher == nil || teacher.ID != discipline.TeacherID {
			return nil, fmt.Errorf("access denied")
		}
	}

	return discipline, nil
}

func (d *DisciplineService) GetDisciplinesByTeacher(ctx context.Context, teacherID int) ([]models.Discipline, error) {
	user, err := d.service.RequireRole(ctx, "teacher", "admin")
	if err != nil {
		return nil, err
	}

	// Преподаватель может видеть только свои дисциплины
	if user.Role == "teacher" {
		teacher, err := d.service.repo.Teachers.GetByID(ctx, user.ID)
		if err != nil || teacher == nil || teacher.ID != teacherID {
			return nil, fmt.Errorf("access denied")
		}
	}

	return d.service.repo.Disciplines.GetByTeacherID(ctx, teacherID)
}

func (d *DisciplineService) GetDisciplinesByGroup(ctx context.Context, groupID int) ([]models.Discipline, error) {
	user, err := d.service.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	// Проверка прав доступа
	switch user.Role {
	case "student":
		student, err := d.service.repo.Students.GetByID(ctx, user.ID)
		if err != nil || student == nil || student.GroupID != groupID {
			return nil, fmt.Errorf("access denied")
		}
	case "teacher":
		// Преподаватель может видеть дисциплины любой группы
	case "admin":
		// Админ может все
	default:
		return nil, fmt.Errorf("access denied")
	}

	return d.service.repo.Disciplines.GetByGroupID(ctx, groupID)
}

func (d *DisciplineService) UpdateDiscipline(ctx context.Context, discipline *models.Discipline) error {
	_, err := d.service.RequireRole(ctx, "teacher", "admin")
	if err != nil {
		return err
	}

	return d.service.repo.Disciplines.Update(ctx, discipline)
}

func (d *DisciplineService) DeleteDiscipline(ctx context.Context, id int) error {
	// Только админ может удалять дисциплины
	_, err := d.service.RequireRole(ctx, "admin")
	if err != nil {
		return err
	}

	return d.service.repo.Disciplines.Delete(ctx, id)
}

func (d *DisciplineService) GetAllDisciplines(ctx context.Context) ([]models.Discipline, error) {
	// Только админ может видеть все дисциплины
	_, err := d.service.RequireRole(ctx, "admin")
	if err != nil {
		return nil, err
	}

	return d.service.repo.Disciplines.GetAll(ctx)
}
