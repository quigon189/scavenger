package dataservice

import (
	"context"
	"fmt"

	"scavenger/core/models"
)

func (s *DataService) CreateLab(ctx context.Context, lab *models.Lab) error {
	// Только преподаватель может создавать лабораторные
	user, err := s.RequireRole(ctx, "teacher")
	if err != nil {
		return err
	}

	// Проверяем, что дисциплина принадлежит преподавателю
	discipline, err := s.dataRepo.Disciplines.GetByID(ctx, lab.DisciplineID)
	if err != nil {
		return fmt.Errorf("failed to get discipline: %v", err)
	}

	teacher, err := s.dataRepo.Teachers.GetByID(ctx, user.ID)
	if err != nil || teacher == nil || discipline.TeacherID != teacher.ID {
		return fmt.Errorf("access denied: discipline does not belong to teacher")
	}

	return s.dataRepo.Labs.Create(ctx, lab)
}

func (s *DataService) GetLab(ctx context.Context, id int) (*models.Lab, error) {
	user, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	lab, err := s.dataRepo.Labs.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Проверка прав доступа
	switch user.Role {
	case "student":
		student, err := s.dataRepo.Students.GetByID(ctx, user.ID)
		if err != nil || student == nil || student.GroupID != lab.Discipline.GroupID {
			return nil, fmt.Errorf("access denied")
		}
	case "teacher":
		teacher, err := s.dataRepo.Teachers.GetByID(ctx, user.ID)
		if err != nil || teacher == nil || teacher.ID != lab.Discipline.TeacherID {
			return nil, fmt.Errorf("access denied")
		}
	}

	// Загружаем файлы лабораторной
	files, err := s.dataRepo.Labs.GetFiles(ctx, id)
	if err == nil {
		lab.Files = files
	}

	return lab, nil
}

func (s *DataService) GetLabsByDiscipline(ctx context.Context, disciplineID int) ([]models.Lab, error) {
	user, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	// Проверяем права доступа к дисциплине
	discipline, err := s.dataRepo.Disciplines.GetByID(ctx, disciplineID)
	if err != nil {
		return nil, err
	}

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

	labs, err := s.dataRepo.Labs.GetByDisciplineID(ctx, disciplineID)
	if err != nil {
		return nil, err
	}

	// Загружаем файлы для каждой лабораторной
	for i := range labs {
		files, err := s.dataRepo.Labs.GetFiles(ctx, labs[i].ID)
		if err == nil {
			labs[i].Files = files
		}
	}

	return labs, nil
}

func (s *DataService) UpdateLab(ctx context.Context, lab *models.Lab) error {
	user, err := s.RequireRole(ctx, "teacher")
	if err != nil {
		return err
	}

	// Проверяем, что лабораторная принадлежит преподавателю
	existingLab, err := s.dataRepo.Labs.GetByID(ctx, lab.ID)
	if err != nil {
		return fmt.Errorf("failed to get lab: %v", err)
	}

	teacher, err := s.dataRepo.Teachers.GetByID(ctx, user.ID)
	if err != nil || teacher == nil || existingLab.Discipline.TeacherID != teacher.ID {
		return fmt.Errorf("access denied: lab does not belong to teacher")
	}

	return s.dataRepo.Labs.Update(ctx, lab)
}

func (s *DataService) DeleteLab(ctx context.Context, id int) error {
	user, err := s.RequireRole(ctx, "teacher")
	if err != nil {
		return err
	}

	// Проверяем, что лабораторная принадлежит преподавателю
	lab, err := s.dataRepo.Labs.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get lab: %v", err)
	}

	teacher, err := s.dataRepo.Teachers.GetByID(ctx, user.ID)
	if err != nil || teacher == nil || lab.Discipline.TeacherID != teacher.ID {
		return fmt.Errorf("access denied: lab does not belong to teacher")
	}

	return s.dataRepo.Labs.Delete(ctx, id)
}

func (s *DataService) AddFileToLab(ctx context.Context, labID, fileID int) error {
	user, err := s.RequireRole(ctx, "teacher")
	if err != nil {
		return err
	}

	// Проверяем, что лабораторная принадлежит преподавателю
	lab, err := s.dataRepo.Labs.GetByID(ctx, labID)
	if err != nil {
		return fmt.Errorf("failed to get lab: %v", err)
	}

	teacher, err := s.dataRepo.Teachers.GetByID(ctx, user.ID)
	if err != nil || teacher == nil || lab.Discipline.TeacherID != teacher.ID {
		return fmt.Errorf("access denied: lab does not belong to teacher")
	}

	return s.dataRepo.Labs.AddFile(ctx, labID, fileID)
}
