package dataservice

import (
	"context"
	"fmt"

	"scavenger/internal/models"
)

func (s *DataService) CreateGroup(ctx context.Context, group *models.Group) error {
	// Только админ может создавать группы
	_, err := s.RequireRole(ctx, "admin")
	if err != nil {
		return err
	}

	return s.dataRepo.Groups.Create(ctx, group)
}

func (s *DataService) GetGroup(ctx context.Context, id int) (*models.Group, error) {
	user, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}
	
	if user.Role == "student" {
		student, err := s.dataRepo.Students.GetByID(ctx, user.ID)
		if err != nil {
			return nil, err
		}

		if student.GroupID != id {
			return nil, fmt.Errorf("access denied")
		}
	}

	group, err := s.dataRepo.Groups.GetWithStudents(ctx, id)
	if err != nil {
		return nil, err
	}

	// Загружаем дисциплины группы
	disciplines, err := s.dataRepo.Disciplines.GetByGroupID(ctx, id)
	if err == nil {
		group.Disciplines = disciplines
	}

	return group, nil
}

func (s *DataService) UpdateGroup(ctx context.Context, group *models.Group) error {
	// Только админ может обновлять группы
	_, err := s.RequireRole(ctx,"admin")
	if err != nil {
		return err
	}

	return s.dataRepo.Groups.Update(ctx, group)
}

func (s *DataService) DeleteGroup(ctx context.Context, id int) error {
	// Только админ может удалять группы
	_, err := s.RequireRole(ctx, "admin")
	if err != nil {
		return err
	}

	return s.dataRepo.Groups.Delete(ctx, id)
}

func (s *DataService) GetAllGroups(ctx context.Context) ([]models.Group, error) {
	return s.dataRepo.Groups.GetAll(ctx)
}
