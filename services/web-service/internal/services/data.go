package services

import (
	"context"
	"web-service/internal/models"
	"web-service/pkg/apiclient"
)

type DataService struct {
	dataClient *apiclient.DataClient
}

func NewDataService(dataClient *apiclient.DataClient) *DataService {
	return &DataService{dataClient: dataClient}
}

func (s *DataService) GetGroups(ctx context.Context) ([]models.Group, error) {
	groups := []models.Group{}

	apiGroups, err := s.dataClient.GetGroups(ctx)
	if err != nil {
		return groups, err
	}

	for _, g := range apiGroups {
		groups = append(groups, models.Group{
			ID:   g.ID,
			Name: g.Name,
		})
	}

	return groups, nil
}
