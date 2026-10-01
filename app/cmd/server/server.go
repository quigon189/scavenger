package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"scavenger/internal/config"
	"scavenger/internal/handler"
	"scavenger/internal/service"
)

func RunServer(ctx context.Context, cfg *config.Config, svc *service.Services) error {
	h := handler.New(handler.Deps{
		Config:   cfg,
		Services: svc,
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
