package main

import (
	"log"
	"scavenger/internal/config"
	"scavenger/internal/repositories/authrepo"
	"scavenger/internal/repositories/sessionrepo"
	"scavenger/internal/services/authservice"
	"scavenger/pkg/postgres"
	"scavenger/pkg/redis"
)

func main() {
	cfg := config.Load()

	redisClient, err := redis.NewRedisClient(cfg.Redis.URL, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Fatalf("Failed get redis client: %v", err)	
	}

	pgClient, err := postgres.NewPostgresPool(cfg.Postgres.GetDatabaseURL())
	if err != nil {
		log.Fatalf("Failed to get postgres client: %v", err)
	}

	authService := authservice.NewAuthService(
		authrepo.NewPgUserRepository(pgClient),
		sessionrepo.NewSessionRepository(redisClient, cfg.SessionTTL),
		cfg.SessionTTL,
	)


}
