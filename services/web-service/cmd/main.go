package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"web-service/internal/config"
	"web-service/internal/handlers"
	"web-service/internal/session"
	"web-service/pkg/apiclient"
	"web-service/pkg/redisclient"

	"github.com/gorilla/sessions"
)

func main() {
	cfg := config.Load()

	redisClient, err := redisclient.NewRedisClient(cfg.Redis.URL, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Fatalf("ERR Failed to create Redis client: %v", err)
	}

	authClient := apiclient.NewAuthClient("http://"+cfg.Api.AuthURL, time.Duration(cfg.Api.Timeout) * time.Second)
	dataClient := apiclient.NewDataClient("http://"+cfg.Api.DataURL, time.Duration(cfg.Api.Timeout) * time.Second)

	sessionService := session.NewSessionService(
		redisClient, 
		authClient,
		sessions.NewCookieStore([]byte("secret-key")),
	)

	router := handlers.NewRouter(
		sessionService,
		authClient,
		dataClient,
	)

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router.Handler(),
	}

	go func() {
		log.Printf("INFO Data service starting on port %s", cfg.Server.Port)
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
