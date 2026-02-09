package session

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"web-service/internal/models"
	"web-service/pkg/apiclient"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/sessions"
)

type SessionService struct {
	sm            *sessions.CookieStore
	redisClient   *redis.Client
	authClient    *apiclient.AuthClient
	sessionPrefix string
	httpClent     *http.Client
}

const sessionKey = "session"

func NewSessionService(redisClient *redis.Client, authClient *apiclient.AuthClient, cookieStore *sessions.CookieStore) *SessionService {
	return &SessionService{
		sm:            cookieStore,
		redisClient:   redisClient,
		authClient:    authClient,
		sessionPrefix: "session:",
		httpClent: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type UserSession struct {
	User      models.User `json:"user"`
	ExpiresAt time.Time   `json:"expires_at"`
}

func (s *SessionService) SetSession(w http.ResponseWriter, r *http.Request, sessionID string, expires_at time.Time) {
	session, err := s.sm.Get(r, sessionKey)
	if err != nil {
		log.Printf("WARN Failed to get session: %v", err)
		return
	}

	session.Values["session_id"] = sessionID
	session.Save(r, w)
}

func (s *SessionService) GetSession(r *http.Request) (*UserSession, string, error) {
	session, err := s.sm.Get(r, sessionKey)
	if err != nil {
		return nil, "", err
	}

	sessionID, ok := session.Values["session_id"].(string)
	if !ok {
		return nil, "", fmt.Errorf("Failed to get session_id from cookie session")
	}

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

	userSession, err := s.getSessionViaAuth(r.Context(), sessionID)
	if err == nil {
		return userSession, sessionID, nil
	}
	return nil, "", err
}

func (s *SessionService) DeleteSession(w http.ResponseWriter, r *http.Request) {
	session, err := s.sm.Get(r, sessionKey)
	if err != nil {
		log.Println("WARN Failed to get session to delete")
	}

	session.Options.MaxAge = -1
	session.Save(r, w)
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
		Status:   authResp.User.Status,
		GroupID:  authResp.User.GroupID,
	}

	return &UserSession{User: user, ExpiresAt: authResp.ExpiresAt}, nil
}
