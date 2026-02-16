package dataservice

import (
	"context"

	"scavenger/internal/models"
)

func (s *DataService) CreatePeriod(ctx context.Context, period *models.Period) error {
	// Только админ может создавать периоды
	_, err := s.RequireRole(ctx, "admin")
	if err != nil {
		return err
	}

	return s.dataRepo.Periods.Create(ctx, period)
}

func (s *DataService) GetPeriod(ctx context.Context, id int) (*models.Period, error) {
	_, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	return s.dataRepo.Periods.GetByID(ctx, id)
}

func (s *DataService) GetActivePeriod(ctx context.Context) (*models.Period, error) {
	_, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	return s.dataRepo.Periods.GetActive(ctx)
}

func (s *DataService) UpdatePeriod(ctx context.Context, period *models.Period) error {
	// Только админ может обновлять периоды
	_, err := s.RequireRole(ctx, "admin")
	if err != nil {
		return err
	}

	return s.dataRepo.Periods.Update(ctx, period)
}

func (s *DataService) DeletePeriod(ctx context.Context, id int) error {
	// Только админ может удалять периоды
	_, err := s.RequireRole(ctx, "admin")
	if err != nil {
		return err
	}

	return s.dataRepo.Periods.Delete(ctx, id)
}

func (s *DataService) GetAllPeriods(ctx context.Context) ([]models.Period, error) {
	_, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	return s.dataRepo.Periods.GetAll(ctx)
}
