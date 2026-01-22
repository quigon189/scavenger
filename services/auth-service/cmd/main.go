package main

import (
	"auth-service/internal/config"
	"auth-service/internal/handlers"
	"auth-service/internal/repository/postgres"
	"auth-service/internal/repository/redis"
	"auth-service/internal/services"
	"auth-service/pkg/database"
	"auth-service/pkg/redisclient"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()

	dbPool, err := database.NewPostgresPool(cfg.GetDatabaseURL())
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer dbPool.Close()

	redisClient, err := redisclient.NewRedisClient(cfg.RedisURL, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	userRepo := postgres.NewUserRepository(dbPool)
	sessionRepo := redis.NewSessionRepository(redisClient, cfg.SessionTTL)

	authService := services.NewAuthService(userRepo, sessionRepo, cfg.SessionTTL)

	authHandler := handlers.NewAuthHandler(authService)

	router := http.NewServeMux()
	authHandler.RegisterRoutes(router)

	server := &http.Server{
		Addr: ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		log.Printf("Auth service start on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stoped")
}
