package server

import (
	"context"
	"log"
	"net/http"
	"scavenger/internal/database"
	filestorage "scavenger/internal/file_storage"
	"scavenger/internal/models"
)

type Server struct {
	httpServer *http.Server
	config *models.Config
	db *database.Database
	fs *filestorage.FileStorage
}

func New(cfg *models.Config, db *database.Database, fs *filestorage.FileStorage) *Server {
	router := setupRouter(cfg, db, fs)
	return &Server{
		httpServer: &http.Server{
			Addr: ":"+cfg.Server.Port,
			Handler: *router,
		},
		config: cfg,
		db: db,
		fs: fs,
	}
}

func (s *Server) Start() error {
	log.Printf("Starting server")
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
