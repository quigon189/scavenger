package session

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"data-service/internal/models"

	"github.com/go-redis/redis/v8"
)

type SessionService struct {
	redisClient    *redis.Client
	authServiceURL string
	sessionPrefix  string
	httpClent      *http.Client
}

func NewSessionService(redisClient *redis.Client, authServiceURL string) *SessionService {
	return &SessionService{
		redisClient:    redisClient,
		authServiceURL: authServiceURL,
		sessionPrefix:  "session:",
		httpClent: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type UserSession struct {
	User      models.User `json:"user"`
	ExpiresAt time.Time   `json:"expires_at"`
}

func (s *SessionService) GetSession(ctx context.Context, sessionID string) (*UserSession, error) {
	key := s.sessionPrefix + sessionID
	data, err := s.redisClient.Get(ctx, key).Bytes()
	if err == nil {
		var session UserSession
		if err := json.Unmarshal(data, &session); err == nil {
			return &session, nil
		}
	}

	return s.getSessionViaAuth(ctx, sessionID)
}

func (s *SessionService) getSessionViaAuth(ctx context.Context, sessionID string) (*UserSession, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", s.authServiceURL+"/validate", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("session_id", sessionID)

	resp, err := s.httpClent.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call auth service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auth sevice returned %d", resp.StatusCode)
	}

	var authResp struct {
		User      models.User `json:"user"`
		SessionID string      `json:"session_id"`
		ExpiresAt time.Time   `json:"expires_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &UserSession{User: authResp.User, ExpiresAt: authResp.ExpiresAt}, nil
}
