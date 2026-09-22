package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"scavenger/core/config"
	"scavenger/core/repositories/authrepo"
	"scavenger/core/repositories/datarepo"
	"scavenger/core/repositories/sessionrepo"
	"scavenger/core/services"
	"scavenger/core/services/authservice"
	"scavenger/core/services/dataservice"
	"scavenger/core/storage"
	"scavenger/web/internal/handlers"
	"scavenger/web/internal/services/session"
	"scavenger/web/pkg/minio"
	"scavenger/web/pkg/postgres"
	"scavenger/web/pkg/redis"
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
	dataService := dataservice.NewDataService(dataRepo, sessionRepo, storageMinio, cfg.Minio.Bucket, time.Duration(cfg.Minio.SignedURLTTL)*time.Second)
	sessionService := session.NewSessionService(authService, cookieStore)

	svcs := services.NewServices(authService, dataService)

	router := handlers.NewRouter(svcs, sessionService)

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
