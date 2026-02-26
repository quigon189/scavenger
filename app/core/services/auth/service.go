package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"scavenger/core/models"
	"scavenger/core/repositories"
	"scavenger/core/utils"

	"github.com/google/uuid"
)

type AuthService struct {
	userRepo    repositories.UserRepo
	sessionRepo repositories.SessionRepo
	sessionTTL  time.Duration
}

func NewAuthService(authRepo repositories.UserRepo, sessionRepo repositories.SessionRepo, sessionTTL time.Duration) *AuthService {
	return &AuthService{
		userRepo:    authRepo,
		sessionRepo: sessionRepo,
		sessionTTL:  sessionTTL,
	}
}

func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.UserSession, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	if user == nil {
		return nil, fmt.Errorf("invalid cerdentials")
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid cerdentials")
	}

	sessionID := uuid.New().String()
	session := &models.UserSession{
		User: *user,
		ID:   sessionID,
	}

	if err := s.sessionRepo.SetSession(ctx, sessionID, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	return session, nil
}

func (s *AuthService) ValidateSession(ctx context.Context, sessionID string) (*models.UserSession, error) {
	session, err := s.sessionRepo.GetSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %v", err)
	}
	if session == nil {
		return nil, fmt.Errorf("session not found or expired")
	}

	if time.Until(session.ExpiresAt) < s.sessionTTL/2 {
		if err := s.sessionRepo.RefreshSession(ctx, sessionID, session); err != nil {
			log.Printf("Failed to refresh session: %v", err)
		}
	}

	return session, nil
}

func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	return s.sessionRepo.DeleteSession(ctx, sessionID)
}

func (s *AuthService) LogoutAll(ctx context.Context, userID int) error {
	sessionIDs, err := s.sessionRepo.GetUserSessions(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user sessions: %v", err)
	}

	for _, sessionID := range sessionIDs {
		s.sessionRepo.DeleteSession(ctx, sessionID)
	}

	return nil
}
