package services

import (
	"context"
	"fmt"
	"log"
	"regexp"
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
			ID:        g.ID,
			Name:      g.Name,
			CreatedAt: g.CreatedAt,
		})
	}

	return groups, nil
}

func (s *DataService) CreateGroup(ctx context.Context, sessionID string, name string) error {
	re := regexp.MustCompile(`^[А-Я]+-\d+$`)

	if !re.MatchString(name) {
		return fmt.Errorf("Недопустимое имя группы")
	}

	req := apiclient.CreateGroupRequest{
		Name: name,
	}

	_, err := s.dataClient.CreateGroup(ctx, sessionID, req)
	if err != nil {
		log.Printf("ERR Failed to create group: %v", err)
		return fmt.Errorf("Ошибка при создании группы: %v", err)
	}

	return nil
}

func (s *DataService) DeleteGroup(ctx context.Context, sessionID string, id int) error {
	err := s.dataClient.DeleteGroup(ctx, sessionID, id)
	if err != nil {
		log.Printf("ERR Failed to delete group: %v", err)
		return fmt.Errorf("Ошибка при удалении группы: %w", err)
	}
	return nil
}

func (s *DataService) GetAdminDashboard(ctx context.Context, sessionID string) (models.AdminDashboardStats, error) {
	// заглушка
	stats := models.AdminDashboardStats{
		TotalStudents:    10,
		TotalTeachers:    2,
		TotalGroups:      2,
		TotalDisciplines: 5,
		TotalLabs:        10,
		PendingStudents:  3,
		RecentReports:    100,
	}

	return stats, nil
}
