package service

import (
	"context"
	"errors"
	"scavenger/internal/auth"
	"scavenger/internal/domain"
	"scavenger/internal/repository"
	"time"
)

var (
	ErrInvalidCredetials = errors.New("invalid credentials")
	ErrUnauthorized      = errors.New("unauthorized")
)

type AuthService struct {
	repo          repository.Repository
	sessionSecret []byte
	sessionTTL    time.Duration
}

func NewAuth(repo repository.Repository, secret []byte) *AuthService {
	return &AuthService{
		repo:          repo,
		sessionSecret: secret,
		sessionTTL:    30 * 24 * time.Hour,
	}
}

type LoginInput struct {
	Email    string
	Password string
}

func (s *AuthService) Login(ctx context.Context, in LoginInput) (string, *domain.User, error) {
	u, err := s.repo.Users().ByEmail(ctx, in.Email)
	if errors.Is(err, repository.ErrNotFound) {
		return "", nil, ErrInvalidCredetials
	}
	if err != nil {
		return "", nil, err
	}
	if !auth.VerifyPassword(in.Password, u.PasswordHash) {
		return "", nil, ErrInvalidCredetials
	}
	sid, err := auth.NewSessionID()
	if err != nil {
		return "", nil, err
	}
	if err := s.repo.Sessions().Create(ctx, &domain.Session{
		ID: sid, UserID: u.ID, ExpiresAt: time.Now().Add(s.sessionTTL),
	}); err != nil {
		return "", nil, err
	}
	return sid, u, nil
}

func (s *AuthService) Logout(ctx context.Context, sid string) error {
	return s.repo.Sessions().Delete(ctx, sid)
}

func (s *AuthService) UserBySession(ctx context.Context, sid string) (*domain.User, error) {
	sess, err := s.repo.Sessions().ByID(ctx, sid)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrUnauthorized
	}
	if err != nil {
		return nil, err
	}
	if time.Now().After(sess.ExpiresAt) {
		_ = s.repo.Sessions().Delete(ctx, sid)
		return nil, ErrUnauthorized
	}
	return s.repo.Users().ByID(ctx, sess.UserID)
}
