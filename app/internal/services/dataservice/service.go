package dataservice

import (
	"context"
	"fmt"
	"time"

	"scavenger/internal/models"
	"scavenger/internal/repositories"
	"scavenger/internal/storage"
)

type DataService struct {
	dataRepo     repositories.DataRepository
	sessionRepo  repositories.SessionRepo
	storage      storage.Storage
	bucketName   string
	signedURLTTL time.Duration
}

func NewDataService(
	dataRepo repositories.DataRepository,
	sessionRepo repositories.SessionRepo,
	storage storage.Storage,
	bucketName string,
	signedURLTTL time.Duration,
) *DataService {
	return &DataService{
		dataRepo:    dataRepo,
		sessionRepo: sessionRepo,
		storage:     storage,
		bucketName:  bucketName,
		signedURLTTL: signedURLTTL,
	}
}
func (s *DataService) Authenticate(ctx context.Context) (*models.User, error) {
	session, ok := ctx.Value("session").(*models.UserSession)
	if !ok {
		return nil, fmt.Errorf("failed to get session")
	}

	return &session.User, nil
}

func (s *DataService) RequireRole(ctx context.Context, requiredRoles ...string) (*models.User, error) {
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
