package dataservice

import (
	"context"
	"scavenger/internal/models"
)

func (s *DataService) GetAdminDashboard(ctx context.Context) (*models.AdminDashboard, error) {
	_, err := s.RequireRole(ctx, "admin")
	if err != nil {
		return nil, err
	}

	return s.dataRepo.Admin.GetAdminDashboard(ctx)
}
