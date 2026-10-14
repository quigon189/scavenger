package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"scavenger/internal/config"
	"scavenger/internal/repository"
	"scavenger/internal/repository/sqlite"
	"scavenger/internal/service"
	"scavenger/internal/storage"
	"scavenger/internal/storage/fs"

	"github.com/lmittmann/tint"
)

func main() {
	if err := run(); err != nil {
		slog.Error("failed", "err", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	setupLogger(cfg.LogLevel)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	repo, err := sqlite.Open(ctx, cfg.DBDSN)
	if err != nil {
		return err
	}
	slog.Debug("repository opened", "dsn", cfg.DBDSN)

	st, err := newStorage(cfg)
	if err != nil {
		return err
	}

	services, err := getServices(cfg, repo, st)
	if err != nil {
		return err
	}

	if len(os.Args) < 2 {
		return RunServer(ctx, cfg, services)
	}

	switch os.Args[1] {
	case "serve":
		return RunServer(ctx, cfg, services)
	case "seed":
		return Seed(ctx, services)
	case "migrate":
		return Migrate(ctx, "sqlite3", repo.DB)
	case "help", "-h", "--help":
		fmt.Println("usage: scavenger [serve|seed]")
	default:
		fmt.Printf("unknown command: %s\n", os.Args[1])
	}

	return nil
}

func setupLogger(level string) {
	var lvl slog.Level
	_ = lvl.UnmarshalText([]byte(level))
	logger := slog.New(tint.NewTextHandler(os.Stdout, &tint.Options{
		Level:      lvl,
		TimeFormat: time.Kitchen,
		NoColor:    false,
	}))
	slog.SetDefault(logger)
}

func getServices(cfg *config.Config, repo repository.Repository, st storage.Storage) (*service.Services, error) {
	services := service.New(service.Deps{
		Repo:          repo,
		Storage:       st,
		MaxUploadSize: cfg.MaxUploadSize,
		SessionSecret: []byte("123"),
	})

	return services, nil
}

func newStorage(cfg *config.Config) (storage.Storage, error) {
	switch cfg.StorageDriver {
	case "fs":
		return fs.New(cfg.StorageFSRoot)
	case "s3":
		return nil, errors.New("s3 storage not implemented yet")
	default:
		return nil, fmt.Errorf("unknown storage driver %s", cfg.StorageDriver)
	}
}
