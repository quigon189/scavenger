package dataservice

import (
	"context"
	"scavenger/core/models"
)

func (s *DataService) GetAdminDashboard(ctx context.Context) (*models.AdminDashboard, error) {
	_, err := s.RequireRole(ctx, "admin")
	if err != nil {
		return nil, err
	}

	return s.dataRepo.Admin.GetAdminDashboard(ctx)
}

func (s *DataService) ConfrimStudent(ctx context.Context, user *models.User) error {
	_, err := s.RequireRole(ctx, "admin")
	if err != nil {
		return err
	}

	return s.dataRepo.Students.Create(ctx, &models.Student{
		ID: user.ID,
		GroupID: *user.GroupID,
	})
}
