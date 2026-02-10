package services

import (
	"context"
	"fmt"
	"net/mail"
	"regexp"
	"slices"
	"strings"
	"web-service/internal/models"
	"web-service/internal/session"
	"web-service/pkg/apiclient"
)

type AuthService struct {
	authClient *apiclient.AuthClient
}

func NewAuthService(authClient *apiclient.AuthClient) *AuthService {
	return &AuthService{authClient: authClient}
}

func (s *AuthService) Login(ctx context.Context, username string, password string) (*session.UserSession, error) {
	req := apiclient.LoginRequest{
		Username: username,
		Password: password,
	}
	resp, err := s.authClient.Login(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.SessionID == "" {
		return nil, fmt.Errorf("failed to login")
	}

	userSession := session.UserSession{
		User: models.User{
			Username: resp.User.Username,
			Email:    resp.User.Email,
			Name:     resp.User.Name,
			Role:     resp.User.Role,
			Status:   resp.User.Status,
			GroupID:  resp.User.GroupID,
		},
		ID:        resp.SessionID,
		ExpiresAt: resp.ExpiresAt,
	}

	return &userSession, nil
}

func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	return s.authClient.Logout(ctx, sessionID)
}

func (s *AuthService) Register(ctx context.Context, username, password, name, email, role string, group_id int) (*models.User, error) {
	reUsername := regexp.MustCompile(`^[a-zA-Z0-9._-]{3,20}$`)
	reName := regexp.MustCompile(`^[А-ЯЁ][а-яё]+ [А-ЯЁ][а-яё]+( [А-ЯЁ][а-яё]+)?$`)

	errors := []string{}
	if !reUsername.MatchString(username) {
		errors = append(errors, "неверно указан логин")
	}
	if len(password) < 8 {
		errors = append(errors, "пароль не соответствует требованиям")
	}
	if !reName.MatchString(name) {
		errors = append(errors, "неверно указано имя")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		errors = append(errors, "неверно указана электронная почта")
	}
	if !slices.Contains([]string{"student", "admin", "teacher"}, role) {
		errors = append(errors, "неверно указана роль")
	}

	if len(errors) > 0 {
		return nil, fmt.Errorf(strings.Join(errors, "; "))
	}

	req := apiclient.RegisterRequest{
		Username: username,
		Name: name,
		Password: password,
		Email: email,
		Role: role,
		GroupID: &group_id,
	}

	resp, err := s.authClient.Register(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ошибка сервиса регистрации: %v", err)
	}

	return &models.User{
		Username: resp.Username,
		Name: resp.Name,
		Email: resp.Email,
		Role: resp.Role,	
		GroupID: resp.GroupID,
		Status: resp.Status,
	}, nil
}
