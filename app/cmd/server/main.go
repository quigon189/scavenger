package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"scavenger/internal/config"
	"scavenger/internal/repository/sqlite"
	"scavenger/internal/service"

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

	services, err := getServices(ctx, cfg)
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

func getServices(ctx context.Context, cfg *config.Config) (*service.Services, error) {

	repo, err := sqlite.Open(ctx, cfg.DBDSN)
	if err != nil {
		return nil, err
	}
	slog.Debug("repository opened", "dsn", cfg.DBDSN)

	services := service.New(service.Deps{
		Repo:          repo,
		SessionSecret: []byte("123"),
	})

	return services, nil
}
