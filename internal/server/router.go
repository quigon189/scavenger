package server

import (
	"net/http"
	"scavenger/internal/handlers"
	"scavenger/internal/models"
)

func setupRouter(cfg *models.Config) *http.ServeMux {
	mux := http.NewServeMux()

	handler := handlers.NewHandler(cfg)

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", handler.AuthMiddleware(handler.Index))
	mux.HandleFunc("/login", handler.Login)

	return mux
}
