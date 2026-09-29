package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"scavenger/internal/config"
	"scavenger/internal/handler"
	"scavenger/internal/repository/sqlite"
	"scavenger/internal/service"
	"syscall"
	"time"

	"github.com/lmittmann/tint"
)

func main() {
	if err := run(); err != nil {
		slog.Error("server failed", "err", err)
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

	services := service.New(service.Deps{
		Repo:          repo,
		SessionSecret: []byte("123"),
	})

	users, err := repo.Users().List(context.Background())
	if len(users) < 1 {
		err := addTestUsers(services)
		if err != nil {
			return err
		}
	}

	h := handler.New(handler.Deps{
		Config:   cfg,
		Services: services,
	})

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           h.Router(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("http server started", "addr", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		slog.Info("shutdown signal received")
	case err := <-errCh:
		return err
	}

	shtdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return srv.Shutdown(shtdownCtx)
}

func setupLogger(level string) {
	var lvl slog.Level
	_ = lvl.UnmarshalText([]byte(level))
	logger := slog.New(tint.NewTextHandler(os.Stdout, &tint.Options{
		Level: lvl,
		TimeFormat: time.Kitchen,
		NoColor: false,
	}))
	slog.SetDefault(logger)
}

func addTestUsers(svc *service.Services) error {
	users := []service.UserInput{
		{
			Email:    "t@edu",
			Password: "123456",
			Role:     "teacher",
			FullName: "Teacher Teacherov",
		},
		{
			Email:    "s@edu",
			Password: "123456",
			Role:     "student",
			FullName: "Student Studentov",
		},
	}

	for _, ui := range users {
		if u, err := svc.User.Create(context.Background(), ui); err != nil {
			return err
		} else {
			slog.Debug("create user", "user", u)
		}
	}

	createdUsers, err := svc.User.List(context.Background())
	if err != nil {
		return nil
	}

	if len(users) > 1 {
		svc.User.Deactivate(context.Background(), createdUsers[1].ID)
		slog.Debug("user deactivated", "uid", createdUsers[1].ID)
	}

	return nil
}
