package services

import (
	"context"

	"data-service/internal/models"
)

type PeriodService struct {
	service *Service
}

func NewPeriodService(service *Service) *PeriodService {
	return &PeriodService{service: service}
}

func (p *PeriodService) CreatePeriod(ctx context.Context, period *models.Period) error {
	// Только админ может создавать периоды
	_, err := p.service.RequireRole(ctx, "admin")
	if err != nil {
		return err
	}

	return p.service.repo.Periods.Create(ctx, period)
}

func (p *PeriodService) GetPeriod(ctx context.Context, id int) (*models.Period, error) {
	_, err := p.service.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	return p.service.repo.Periods.GetByID(ctx, id)
}

func (p *PeriodService) GetActivePeriod(ctx context.Context) (*models.Period, error) {
	_, err := p.service.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	return p.service.repo.Periods.GetActive(ctx)
}

func (p *PeriodService) UpdatePeriod(ctx context.Context, period *models.Period) error {
	// Только админ может обновлять периоды
	_, err := p.service.RequireRole(ctx, "admin")
	if err != nil {
		return err
	}

	return p.service.repo.Periods.Update(ctx, period)
}

func (p *PeriodService) DeletePeriod(ctx context.Context, id int) error {
	// Только админ может удалять периоды
	_, err := p.service.RequireRole(ctx, "admin")
	if err != nil {
		return err
	}

	return p.service.repo.Periods.Delete(ctx, id)
}

func (p *PeriodService) GetAllPeriods(ctx context.Context) ([]models.Period, error) {
	_, err := p.service.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	return p.service.repo.Periods.GetAll(ctx)
}
