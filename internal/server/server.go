package server

import (
	"context"
	"log"
	"net/http"
	"scavenger/internal/models"
)

type Server struct {
	httpServer *http.Server
	config *models.Config
}

func New(cfg *models.Config) *Server {
	router := setupRouter(cfg)
	return &Server{
		httpServer: &http.Server{
			Addr: ":"+cfg.Server.Port,
			Handler: router,
		},
	}
}

func (s *Server) Start() error {
	log.Printf("Starting server")
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
