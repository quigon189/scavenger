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

	services := service.New(service.Deps{
		Repo: repo,
		SessionSecret: []byte("123"),
	})

	h := handler.New(handler.Deps{
		Config: cfg,
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
	case err := <- errCh:
		return err
	}

	shtdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return srv.Shutdown(shtdownCtx)
}

func setupLogger(level string) {
	var lvl slog.Level
	_ = lvl.UnmarshalText([]byte(level))
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})))
}
