package dataservice

import (
	"context"
	"fmt"

	"scavenger/internal/models"
)

func (s *DataService) CreateDiscipline(ctx context.Context, discipline *models.Discipline) error {
	_, err := s.RequireRole(ctx, "teacher", "admin")
	if err != nil {
		return err
	}

	return s.dataRepo.Disciplines.Create(ctx, discipline)
}

func (s *DataService) GetDiscipline(ctx context.Context, id int) (*models.Discipline, error) {
	user, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	discipline, err := s.dataRepo.Disciplines.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Проверка прав доступа
	switch user.Role {
	case "student":
		student, err := s.dataRepo.Students.GetByID(ctx, user.ID)
		if err != nil || student == nil || student.GroupID != discipline.GroupID {
			return nil, fmt.Errorf("access denied")
		}
	case "teacher":
		teacher, err := s.dataRepo.Teachers.GetByID(ctx, user.ID)
		if err != nil || teacher == nil || teacher.ID != discipline.TeacherID {
			return nil, fmt.Errorf("access denied")
		}
	}

	return discipline, nil
}

func (s *DataService) GetDisciplinesByTeacher(ctx context.Context, teacherID int) ([]models.Discipline, error) {
	user, err := s.RequireRole(ctx, "teacher", "admin")
	if err != nil {
		return nil, err
	}

	// Преподаватель может видеть только свои дисциплины
	if user.Role == "teacher" {
		teacher, err := s.dataRepo.Teachers.GetByID(ctx, user.ID)
		if err != nil || teacher == nil || teacher.ID != teacherID {
			return nil, fmt.Errorf("access denied")
		}
	}

	return s.dataRepo.Disciplines.GetByTeacherID(ctx, teacherID)
}

func (d *DataService) GetDisciplinesByGroup(ctx context.Context, groupID int) ([]models.Discipline, error) {
	user, err := d.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	// Проверка прав доступа
	switch user.Role {
	case "student":
		student, err := d.dataRepo.Students.GetByID(ctx, user.ID)
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

	return d.dataRepo.Disciplines.GetByGroupID(ctx, groupID)
}

func (d *DataService) UpdateDiscipline(ctx context.Context, discipline *models.Discipline) error {
	_, err := d.RequireRole(ctx, "teacher", "admin")
	if err != nil {
		return err
	}

	return d.dataRepo.Disciplines.Update(ctx, discipline)
}

func (d *DataService) DeleteDiscipline(ctx context.Context, id int) error {
	// Только админ может удалять дисциплины
	_, err := d.RequireRole(ctx, "admin")
	if err != nil {
		return err
	}

	return d.dataRepo.Disciplines.Delete(ctx, id)
}

func (d *DataService) GetAllDisciplines(ctx context.Context) ([]models.Discipline, error) {
	// Только админ может видеть все дисциплины
	_, err := d.RequireRole(ctx, "admin")
	if err != nil {
		return nil, err
	}

	return d.dataRepo.Disciplines.GetAll(ctx)
}
