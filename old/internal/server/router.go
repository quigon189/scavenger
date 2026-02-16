package server

import (
	"net/http"
	"scavenger/internal/database"
	filestorage "scavenger/internal/file_storage"
	"scavenger/internal/handlers"
	"scavenger/internal/models"
)

func setupRouter(cfg *models.Config, db *database.Database, fileStorage *filestorage.FileStorage) *http.Handler {
	mux := http.NewServeMux()

	handler := handlers.NewHandler(cfg, db, fileStorage)

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))


	mux.HandleFunc("/404", handler.NotFound)
	mux.HandleFunc("/", handler.AuthMiddleware(handler.Index))
	mux.HandleFunc("/register", handler.StudentRegisterPage)
	mux.HandleFunc("POST /register", handler.RegisterStudent)
	mux.HandleFunc("/dashboard", handler.AuthMiddleware(handler.Dashboard))
	mux.HandleFunc("/logout", handler.AuthMiddleware(handler.Logout))
	mux.HandleFunc("POST /change-password", handler.AuthMiddleware(handler.ChangePassword))
	mux.HandleFunc("/change-theme", handler.AuthMiddleware(handler.ChangeTheme))

	mux.HandleFunc("/disciplines/{id}", handler.StudentMiddleware(handler.DisciplinePage))
	mux.HandleFunc("/disciplines/{discID}/labs/{labID}", handler.AuthMiddleware(handler.LabMarkdownPage))
	mux.HandleFunc("GET /disciplines/{discID}/labs/{labID}/reports", handler.StudentMiddleware(handler.LabReportPage))
	mux.HandleFunc("POST /disciplines/{discID}/labs/{labID}/reports", handler.StudentMiddleware(handler.UploadLabReport))

	mux.HandleFunc("GET /admin/groups", handler.AdminMiddleware(handler.GroupManager))
	mux.HandleFunc("POST /admin/groups", handler.AdminMiddleware(handler.AddGroup))
	mux.HandleFunc("POST /admin/groups/{id}", handler.AdminMiddleware(handler.EditGroup))
	mux.HandleFunc("POST /admin/groups/{groupID}/delete", handler.AdminMiddleware(handler.DeleteGroup))
	mux.HandleFunc("POST /admin/groups/{id}/disciplines", handler.AdminMiddleware(handler.AddDiscToGroup))
	mux.HandleFunc("POST /admin/groups/{groupID}/disciplines/{discID}/remove", handler.AdminMiddleware(handler.RemoveDiscFromGroup))

	mux.HandleFunc("GET /admin/disciplines", handler.AdminMiddleware(handler.DisciplinesManager))
	mux.HandleFunc("POST /admin/disciplines", handler.AdminMiddleware(handler.AddDiscipline))
	mux.HandleFunc("POST /admin/disciplines/{id}", handler.AdminMiddleware(handler.EditDiscipline))
	mux.HandleFunc("GET /admin/disciplines/{id}", handler.AdminMiddleware(handler.DisciplineLabs))
	mux.HandleFunc("POST /admin/disciplines/{id}/delete", handler.AdminMiddleware(handler.DeleteDiscipline))

	mux.HandleFunc("POST /admin/disciplines/{id}/labs", handler.AdminMiddleware(handler.AddDisciplineLabs))
	mux.HandleFunc("POST /admin/disciplines/{discID}/labs/{labID}", handler.AdminMiddleware(handler.EditDisciplineLab))
	mux.HandleFunc("POST /admin/disciplines/{discID}/labs/{labID}/delete", handler.AdminMiddleware(handler.DeleteDisciplineLab))

	mux.HandleFunc("GET /admin/students", handler.AdminMiddleware(handler.StudentsManager))
	mux.HandleFunc("POST /admin/students", handler.AdminMiddleware(handler.AddStudents))
	mux.HandleFunc("POST /admin/students/{username}", handler.AdminMiddleware(handler.EditStudent))
	mux.HandleFunc("POST /admin/students/{username}/delete", handler.AdminMiddleware(handler.DeleteStudent))

	mux.HandleFunc("/admin/reports", handler.AdminMiddleware(handler.ReportsPage))
	mux.HandleFunc("/admin/reports/table", handler.AdminMiddleware(handler.ReportsTable))
	mux.HandleFunc("/admin/reports/labs", handler.AdminMiddleware(handler.ReportsLabsByDiscipline))
	mux.HandleFunc("/admin/reports/{id}/review", handler.AdminMiddleware(handler.ReportReviewPage))
	mux.HandleFunc("POST /admin/reports/{id}/grade", handler.AdminMiddleware(handler.GradeReport))

	mux.HandleFunc("/files/", handler.AuthMiddleware(handler.GetFile))

	h := handler.AlertMiddleware(mux)

	return &h
}
