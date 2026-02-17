package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"scavenger/internal/config"
	"scavenger/internal/handlers"
	"scavenger/internal/repositories/authrepo"
	"scavenger/internal/repositories/datarepo"
	"scavenger/internal/repositories/sessionrepo"
	"scavenger/internal/services"
	"scavenger/internal/services/authservice"
	"scavenger/internal/services/dataservice"
	"scavenger/internal/services/sessionservice"
	"scavenger/internal/storage"
	"scavenger/pkg/minio"
	"scavenger/pkg/postgres"
	"scavenger/pkg/redis"
	"syscall"
	"time"

	"github.com/gorilla/sessions"
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

	minioClient, err := minio.NewMinioClient(cfg.Minio.Endpoint, cfg.Minio.AccessKey, cfg.Minio.SecretKey, cfg.Minio.Bucket, cfg.Minio.UseSSL)
	if err != nil {
		log.Fatalf("Failed to get MinIO client: %v", err)
	}

	storageMinio, err := storage.NewMinioStorage(minioClient, cfg.Minio.Bucket)
	if err != nil {
		log.Fatalf("Failed to get MinIO storage: %v", err)
	}

	cookieStore := sessions.NewCookieStore([]byte(cfg.CookieSecret))

	authRepo := authrepo.NewPgUserRepository(pgClient)
	dataRepo := datarepo.NewDataRepository(pgClient)
	sessionRepo := sessionrepo.NewSessionRepository(redisClient, cfg.SessionTTL)

	authService := authservice.NewAuthService(authRepo, sessionRepo, cfg.SessionTTL)
	dataService := dataservice.NewDataService(dataRepo, sessionRepo, storageMinio, cfg.Minio.Bucket, time.Duration(cfg.Minio.SignedURLTTL) * time.Second)
	sessionService := sessionservice.NewSessionService(sessionRepo, cookieStore)

	svcs := services.NewServices(authService, dataService, sessionService)

	router := handlers.NewRouter(svcs)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router.Handler(),
	}

	go func() {
		log.Printf("INFO Data service starting on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ERR Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("INFO Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("ERR Server forced to shutdown: %v", err)
	}

	log.Println("INFO Server stopped")
}
