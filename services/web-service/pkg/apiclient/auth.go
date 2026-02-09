package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AuthClient struct {
	client *Client
}

func NewAuthClient(baseURL string, timeout time.Duration) *AuthClient {
	return &AuthClient{
		client: NewClient(baseURL, timeout),
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"`
	GroupID  *int   `json:"group_id,omitempty"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	GroupID  *int   `json:"group_id"`
	Status   string `json:"status"`
}

type AuthResponse struct {
	User      User      `json:"user"`
	SessionID string    `json:"session_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (a *AuthClient) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	var resp AuthResponse

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	err = a.client.doRequest(ctx, http.MethodPost, "/login", nil, bytes.NewReader(body), &resp)
	if err != nil {
		return nil, fmt.Errorf("login failed: %w", err)
	}

	return &resp, nil
}

func (a *AuthClient) Register(ctx context.Context, req RegisterRequest) (*User, error) {
	var user User

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	err = a.client.doRequest(ctx, http.MethodPost, "/register", nil, bytes.NewReader(body), &user)
	if err != nil {
		return nil, fmt.Errorf("register failed: %w", err)
	}

	return &user, nil
}

func (a *AuthClient) ValidateSession(ctx context.Context, sessionID string) (*AuthResponse, error) {
	var resp AuthResponse

	headers := map[string]string{
		"X-Session-ID": sessionID,
	}

	err := a.client.doRequest(ctx, http.MethodPost, "/validate", headers, nil, &resp)
	if err != nil {
		return nil, fmt.Errorf("session validation failed: %w", err)
	}

	return &resp, nil
}

func (a *AuthClient) Logout(ctx context.Context, sessionID string) error {
	headers := map[string]string{
		"X-Session-ID": sessionID,
	}

	err := a.client.doRequest(ctx, http.MethodGet, "/logout", headers, nil, nil)
	if err != nil {
		return fmt.Errorf("logout failed: %w", err)
	}

	return nil
}

func (a *AuthClient) ChangePassword(ctx context.Context, sessionID, oldPassword, newPassword string) error {
	headers := map[string]string{
		"X-Session-ID": sessionID,
	}

	req := map[string]string{
		"old_password": oldPassword,
		"new_password": newPassword,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	err = a.client.doRequest(ctx, http.MethodPost, "/change_password", headers, bytes.NewReader(body), nil)
	if err != nil {
		return fmt.Errorf("change password failed: %w", err)
	}

	return nil
}
