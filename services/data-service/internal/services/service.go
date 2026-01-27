package services

import (
	"context"
	"fmt"
	"time"

	"data-service/internal/models"
	"data-service/internal/repository/postgres"
	"data-service/internal/services/session"

	"github.com/minio/minio-go/v7"
)

type Service struct {
	repo         *postgres.Repository
	session      *session.SessionService
	minioClient  *minio.Client
	bucketName   string
}

func NewService(
	repo *postgres.Repository,
	sessionService *session.SessionService,
	minioClient *minio.Client,
	bucketName string,
) *Service {
	return &Service{
		repo:        repo,
		session:     sessionService,
		minioClient: minioClient,
		bucketName:  bucketName,
	}
}

func (s *Service) Authenticate(ctx context.Context, sessionID string) (*models.User, error) {
	session, err := s.session.GetSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %v", err)
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, fmt.Errorf("session expired")
	}

	return &session.User, nil
}

func (s *Service) RequireRole(ctx context.Context, sessionID string, requiredRole string) (*models.User, error) {
	user, err := s.Authenticate(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if user.Role != requiredRole && user.Role != "admin" {
		return nil, fmt.Errorf("access denied: required role %s", requiredRole)
	}

	return user, nil
}
