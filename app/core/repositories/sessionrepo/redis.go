package sessionrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"scavenger/core/models"

	"github.com/go-redis/redis/v8"
)

type SessionRepository struct {
	client *redis.Client
	prefix string
	ttl    time.Duration
}

func NewSessionRepository(client *redis.Client, ttl time.Duration) *SessionRepository {
	return &SessionRepository{
		client: client,
		prefix: "session:",
		ttl:    ttl,
	}
}

func (r *SessionRepository) SetSession(ctx context.Context, sessionID string, session *models.UserSession) error {
	session.ExpiresAt = time.Now().Add(r.ttl)

	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %v", err)
	}

	key := r.prefix + sessionID
	err = r.client.Set(ctx, key, data, r.ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to save session: %v", err)
	}

	return nil
}

func (r *SessionRepository) GetSession(ctx context.Context, sessionID string) (*models.UserSession, error) {
	key := r.prefix + sessionID

	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get session: %v", err)
	}

	var session models.UserSession
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %v", err)
	}

	return &session, nil
}

func (r *SessionRepository) DeleteSession(ctx context.Context, sessionID string) error {
	key := r.prefix + sessionID
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete session: %v", err)
	}
	return nil
}

func (r *SessionRepository) RefreshSession(ctx context.Context, sessionID string, session *models.UserSession) error {
	key := r.prefix + sessionID

	session.ExpiresAt = time.Now().Add(r.ttl)

	err := r.client.Expire(ctx, key, r.ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to refresh session: %v", err)
	}

	return nil
}

func (r *SessionRepository) GetUserSessions(ctx context.Context, userID int) ([]string, error) {
	pattern := r.prefix + "*"

	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get session keys: %v", err)
	}

	var userSessionIDs []string
	for _, key := range keys {
		data, err := r.client.Get(ctx, key).Bytes()
		if err != nil {
			continue
		}

		var session models.UserSession
		if err := json.Unmarshal(data, &session); err != nil {
			continue
		}

		if session.User.ID == userID {
			sessionID := key[len(r.prefix):]
			userSessionIDs = append(userSessionIDs, sessionID)
		}
	}

	return userSessionIDs, nil
}

func (r *SessionRepository) HealtCheck(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
