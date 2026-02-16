package authservice

import (
	"context"
	"fmt"
	"log"
	"time"

	"scavenger/internal/models"
	"scavenger/internal/repositories"

	"github.com/google/uuid"
)

type AuthService struct {
	userRepo    repositories.UserRepo
	sessionRepo repositories.SessionRepo
	sessionTTL  time.Duration
}

func NewAuthService(userRepo repositories.UserRepo, sessionRepo repositories.SessionRepo, ttl time.Duration) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		sessionTTL:  ttl,
	}
}

func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
	existingUser, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %v", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("user with username %s already exists", req.Username)
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	user := &models.User{
		Username:     req.Username,
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         req.Role,
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.UserSession, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	if user == nil {
		return nil, fmt.Errorf("invalid cerdentials")
	}

	if !CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid cerdentials")
	}

	sessionID := uuid.New().String()
	session := &models.UserSession{
		User: *user,
		ID: sessionID,
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

func (s *AuthService) ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	if !CheckPasswordHash(oldPassword, user.PasswordHash) {
		return fmt.Errorf("invalid old password")
	}

	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %v", err)
	}

	if err := s.userRepo.UpdatePassword(ctx, user.ID, hashedPassword); err != nil {
		return fmt.Errorf("failed to update password: %v", err)
	}

	return s.LogoutAll(ctx, user.ID)
}
