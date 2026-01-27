package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"data-service/internal/config"
	"data-service/pkg/database"
	"data-service/pkg/minio"
	"data-service/pkg/redisclient"
)

func main() {
	cfg := config.Load()

	_, err := redisclient.NewRedisClient(cfg.RedisURL, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("ERR Failed to create Redis client: %v", err)
	}

	_, err = database.NewPostgresPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("ERR Failed to create Postgres Pool: %v", err)
	}

	_, err = minio.NewMinioClient(cfg.MinioEndpoint, cfg.MinioAccessKey, cfg.MinioSecretKey, cfg.MinioBucket, cfg.MinioUseSSL)
	if err != nil {
		log.Fatalf("ERR Failed to create MinIO client: %v", err)
	}

	router := http.NewServeMux()

	// TODO: регистрация маршрутов

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
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
