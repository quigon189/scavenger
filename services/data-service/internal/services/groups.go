package services

import (
	"context"
	"fmt"

	"data-service/internal/models"
)

type GroupService struct {
	service *Service
}

func NewGroupService(service *Service) *GroupService {
	return &GroupService{service: service}
}

func (g *GroupService) CreateGroup(ctx context.Context, group *models.Group) error {
	// Только админ может создавать группы
	_, err := g.service.RequireRole(ctx, "admin")
	if err != nil {
		return err
	}

	return g.service.repo.Groups.Create(ctx, group)
}

func (g *GroupService) GetGroup(ctx context.Context, id int) (*models.Group, error) {
	user, err := g.service.Authenticate(ctx)
	if err != nil {
		return nil, err
	}
	
	if user.Role == "student" {
		student, err := g.service.repo.Students.GetByID(ctx, user.ID)
		if err != nil {
			return nil, err
		}

		if student.GroupID != id {
			return nil, fmt.Errorf("access denied")
		}
	}

	group, err := g.service.repo.Groups.GetWithStudents(ctx, id)
	if err != nil {
		return nil, err
	}

	// Загружаем дисциплины группы
	disciplines, err := g.service.repo.Disciplines.GetByGroupID(ctx, id)
	if err == nil {
		group.Disciplines = disciplines
	}

	return group, nil
}

func (g *GroupService) UpdateGroup(ctx context.Context, group *models.Group) error {
	// Только админ может обновлять группы
	_, err := g.service.RequireRole(ctx,"admin")
	if err != nil {
		return err
	}

	return g.service.repo.Groups.Update(ctx, group)
}

func (g *GroupService) DeleteGroup(ctx context.Context, id int) error {
	// Только админ может удалять группы
	_, err := g.service.RequireRole(ctx, "admin")
	if err != nil {
		return err
	}

	return g.service.repo.Groups.Delete(ctx, id)
}

func (g *GroupService) GetAllGroups(ctx context.Context) ([]models.Group, error) {
	return g.service.repo.Groups.GetAll(ctx)
}
