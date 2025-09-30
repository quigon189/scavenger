package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"scavenger/internal/config"
	"scavenger/internal/database"
	"scavenger/internal/server"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed load config: %v", err)
	}

	db, err := database.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed connect to db: %v", err)
	}

	err = db.Migrate()
	if err != nil {
		log.Fatalf("Failed apply migrations: %v", err)
	}

	db.SetTestData(cfg)

	srv := server.New(cfg, db)

	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown failed: %v", err)
	}

	log.Println("Server stoped")
}
