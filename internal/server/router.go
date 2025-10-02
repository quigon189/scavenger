package server

import (
	"net/http"
	"scavenger/internal/database"
	"scavenger/internal/handlers"
	"scavenger/internal/models"
)

func setupRouter(cfg *models.Config, db *database.Database) *http.Handler {
	mux := http.NewServeMux()

	handler := handlers.NewHandler(cfg, db)

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))


	mux.HandleFunc("/404", handler.NotFound)
	mux.HandleFunc("/", handler.AuthMiddleware(handler.Index))
	mux.HandleFunc("/dashboard", handler.AuthMiddleware(handler.Dashboard))
	mux.HandleFunc("/logout", handler.AuthMiddleware(handler.Logout))
	mux.HandleFunc("/discipline/{id}", handler.StudentMiddleware(handler.DisciplinePage))

	mux.HandleFunc("/admin/groups", handler.AdminMiddleware(handler.GroupManager))

	mux.HandleFunc("/download/{path...}", handler.AuthMiddleware(handler.DownloadLabs))
	mux.HandleFunc("/upload-report", handler.AuthMiddleware(handler.UploadReport))

	h := handler.AlertMiddleware(mux)

	return &h
}
