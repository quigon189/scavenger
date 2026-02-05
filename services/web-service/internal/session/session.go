package session

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"web-service/internal/models"
	"web-service/pkg/apiclient"

	"github.com/go-redis/redis/v8"
)

type SessionService struct {
	redisClient   *redis.Client
	authClient    *apiclient.AuthClient
	sessionPrefix string
	httpClent     *http.Client
}

func NewSessionService(redisClient *redis.Client, authClient *apiclient.AuthClient) *SessionService {
	return &SessionService{
		redisClient:   redisClient,
		authClient:    authClient,
		sessionPrefix: "session:",
		httpClent: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type UserSession struct {
	User      models.User    `json:"user"`
	ExpiresAt time.Time      `json:"expires_at"`
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

	sessionID := sessionCookie.Value

	key := s.sessionPrefix + sessionID
	data, err := s.redisClient.Get(r.Context(), key).Bytes()
	if err == nil {
		var session UserSession
		if err := json.Unmarshal(data, &session); err == nil {
			return &session, sessionID, nil
		}
	} else {
		log.Printf("WARN Failed to get session from Redis")
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
	authResp, err := s.authClient.ValidateSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Name:     authResp.User.Name,
		Username: authResp.User.Username,
		Role:     authResp.User.Role,
		Email:    authResp.User.Email,
	}

	return &UserSession{User: user, ExpiresAt: authResp.ExpiresAt}, nil
}
