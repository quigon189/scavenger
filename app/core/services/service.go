package services

import (
	"context"
	"fmt"
	"scavenger/core/models"
)

type BaseService struct {
}

func (s *BaseService) Authenticate(ctx context.Context) (*models.User, error) {
	session, ok := ctx.Value("session").(*models.UserSession)
	if !ok {
		return nil, fmt.Errorf("failed to get session")
	}

	return &session.User, nil
}

func (s *BaseService) RequireRole(ctx context.Context, requiredRoles ...string) error {
	user, err := s.Authenticate(ctx)
	if err != nil {
		return err
	}

	access := false
	for _, requiredRole := range requiredRoles {
		if user.Role == requiredRole {
			access = true
			break
		}
	}

	if access {
		return nil
	}
	return fmt.Errorf("access denied")
}
