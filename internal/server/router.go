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

	mux.HandleFunc("GET /admin/groups", handler.AdminMiddleware(handler.GroupManager))
	mux.HandleFunc("POST /admin/groups", handler.AdminMiddleware(handler.AddGroup))
	mux.HandleFunc("POST /admin/groups/{groupID}/delete", handler.AdminMiddleware(handler.DeleteGroup))

	mux.HandleFunc("/admin/disciplines", handler.AdminMiddleware(handler.DisciplinesManager))
	mux.HandleFunc("/admin/students", handler.AdminMiddleware(handler.StudentsManager))

	mux.HandleFunc("/download/{path...}", handler.AuthMiddleware(handler.DownloadLabs))
	mux.HandleFunc("/upload-report", handler.AuthMiddleware(handler.UploadReport))

	h := handler.AlertMiddleware(mux)

	return &h
}
