package services

import (
	"context"
	"fmt"

	"data-service/internal/models"
	"data-service/internal/repository/postgres"
	"data-service/internal/session"
	"data-service/internal/storage"
)

type Service struct {
	repo       *postgres.Repository
	session    *session.SessionService
	storage    storage.Storage
	bucketName string
}

func NewService(
	repo *postgres.Repository,
	sessionService *session.SessionService,
	storage storage.Storage,
	bucketName string,
) *Service {
	return &Service{
		repo:       repo,
		session:    sessionService,
		storage:    storage,
		bucketName: bucketName,
	}
}
func (s *Service) Authenticate(ctx context.Context) (*models.User, error) {
	session, ok := ctx.Value("session").(*session.UserSession)
	if !ok {
		return nil, fmt.Errorf("failed to get session")
	}

	return &session.User, nil
}

func (s *Service) RequireRole(ctx context.Context, requiredRoles ...string) (*models.User, error) {
	user, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	access := false
	for _, requiredRole := range requiredRoles {
		if user.Role == requiredRole {
			access = true
			break
		}
	}

	if access {
		return user, nil
	}
	return nil, fmt.Errorf("access denied")
}
