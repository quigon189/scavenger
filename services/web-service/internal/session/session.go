package session

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

type UserSession struct {
	User      User      `json:"user"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (s *SessionService) SetSession(w http.ResponseWriter, sessionID string, expires_at time.Time) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  expires_at,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}

	http.SetCookie(w, cookie)
}

func (s *SessionService) GetSession(r *http.Request) (*UserSession, string, error) {
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, "", err
	}

	sessionID := sessionCookie.String()

	key := s.sessionPrefix + sessionID
	data, err := s.redisClient.Get(r.Context(), key).Bytes()
	if err == nil {
		var session UserSession
		if err := json.Unmarshal(data, &session); err == nil {
			return &session, sessionID, nil
		}
	}

	session, err := s.getSessionViaAuth(r.Context(), sessionID)
	if err == nil {
		return session, sessionID, nil
	}
	return nil, "", err
}

func (s *SessionService) DeleteSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
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
		User      User      `json:"user"`
		SessionID string    `json:"session_id"`
		ExpiresAt time.Time `json:"expires_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &UserSession{User: authResp.User, ExpiresAt: authResp.ExpiresAt}, nil
}
