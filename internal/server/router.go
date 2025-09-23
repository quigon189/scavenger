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


	mux.HandleFunc("/404", handler.NotFound)
	mux.HandleFunc("/", handler.AuthMiddleware(handler.Index))
	mux.HandleFunc("/dashboard", handler.AuthMiddleware(handler.Dashboard))
	mux.HandleFunc("/logout", handler.AuthMiddleware(handler.Logout))
	mux.HandleFunc("/discipline/{id}", handler.StudentMiddleware(handler.DisciplinePage))

	mux.HandleFunc("/download/{path}", handler.AuthMiddleware(handler.DownloadLabs))

	return mux
}
