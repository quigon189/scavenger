package main

import (
	"log"
	"scavenger/internal/config"
	"scavenger/internal/repositories/authrepo"
	"scavenger/internal/repositories/datarepo"
	"scavenger/internal/repositories/sessionrepo"
	"scavenger/internal/services/authservice"
	"scavenger/internal/services/dataservice"
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

	authRepo := authrepo.NewPgUserRepository(pgClient)
	dataRepo := datarepo.NewDataRepository(pgClient)
	sessionRepo := sessionrepo.NewSessionRepository(redisClient, cfg.SessionTTL)

	authService := authservice.NewAuthService(authRepo, sessionRepo, cfg.SessionTTL)
	dataService := dataservice.NewDataService(*dataRepo)
}
