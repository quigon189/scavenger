package redis

import (
	"auth-service/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type SessionRepository struct {
	client *redis.Client
	ttl    time.Duration
}

func NewSessionRepository(client *redis.Client, ttl time.Duration) *SessionRepository {
	return &SessionRepository{
		client: client,
		ttl:    ttl,
	}
}

func (r *SessionRepository) SetSession(ctx context.Context, token string, user *models.User, expiresAt time.Time) error {
	userJSON, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %v", err)
	}

	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return fmt.Errorf("session already expired")
	}

	return r.client.Set(ctx, fmt.Sprintf("session:%s", token), userJSON, ttl).Err()
}

func (r *SessionRepository) GetSession(ctx context.Context, token string) (*models.User, error) {
	key := fmt.Sprintf("session:%s", token)
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get session: %v", err)
	}

	var user models.User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %v", err)
	}

	return &user, nil
}

func (r *SessionRepository) DeleteSession(ctx context.Context, token string) error {
	key := fmt.Sprintf("session:%s", token)
	return r.client.Del(ctx, key).Err()
}

func (r *SessionRepository) RefreshSession(ctx context.Context, token string, expiresAt time.Time) error {
	key := fmt.Sprintf("session:%s", token)
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return fmt.Errorf("session already expired")
	}

	return r.client.Expire(ctx, key, ttl).Err()
}

func (r *SessionRepository) HealtCheck(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
