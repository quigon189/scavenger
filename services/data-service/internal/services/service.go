package services

import (
	"context"
	"fmt"

	"data-service/internal/models"
	"data-service/internal/repository/postgres"
	"data-service/internal/session"

	"github.com/minio/minio-go/v7"
)

type Service struct {
	repo        *postgres.Repository
	session     *session.SessionService
	minioClient *minio.Client
	bucketName  string
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
func (s *Service) Authenticate(ctx context.Context) (*models.User, error) {
	session, ok := ctx.Value("session").(session.UserSession)
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

	if user.Role == "admin" {
		return user, nil
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
