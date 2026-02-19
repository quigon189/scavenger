package handlers

import (
	"net/http"
	"scavenger/web/internal/middlewares"
	"scavenger/core/models"
	"scavenger/core/services"
)

type Router struct {
	services   *services.Services
	middleware *middlewares.Middleware
	handler    http.Handler
	mux        *http.ServeMux
}

func NewRouter(services *services.Services) *Router {
	router := &Router{
		mux:        http.NewServeMux(),
		services:   services,
		middleware: middlewares.NewMiddleware(services.Session),
	}

	router.registerRoutes()

	router.handler = router.mux

	return router
}

func getSession(r *http.Request) *models.UserSession {
	return r.Context().Value("session").(*models.UserSession)
}

func (r *Router) registerRoutes() {
	authHandler := NewAuthHandler(r.services)
	adminHandler := NewAdminHandler(r.services)
	teacherHandler := NewTeacherHandler(r.services)

	r.mux.HandleFunc("/", r.middleware.SessionRequire(Home))
	r.mux.HandleFunc("/login", authHandler.Login)
	r.mux.HandleFunc("POST /logout", r.middleware.SessionRequire(authHandler.Logout))

	r.mux.HandleFunc("GET /register", authHandler.EnterCode)
	r.mux.HandleFunc("POST /register", authHandler.VerifyCode)
	r.mux.HandleFunc("GET /register/complete", authHandler.RegisterComplete)
	r.mux.HandleFunc("POST /register/complete", authHandler.RegisterCompletePost)

	r.mux.HandleFunc("/admin/dashboard", r.middleware.AdminRequire(adminHandler.Dashboard))

	r.mux.HandleFunc("GET /admin/groups", r.middleware.AdminRequire(adminHandler.Groups))
	r.mux.HandleFunc("POST /admin/groups/{id}/delete", r.middleware.AdminRequire(adminHandler.DeleteGroup))
	r.mux.HandleFunc("POST /admin/groups/create", r.middleware.AdminRequire(adminHandler.CreateGroup))

	r.mux.HandleFunc("GET /admin/teachers", r.middleware.AdminRequire(adminHandler.Teachers))
	r.mux.HandleFunc("POST /admin/teachers/{id}/delete", r.middleware.AdminRequire(adminHandler.DeleteTeacher))
	r.mux.HandleFunc("POST /admin/teachers/create", r.middleware.AdminRequire(adminHandler.CreateTeacher))

	r.mux.HandleFunc("GET /admin/students", r.middleware.AdminRequire(adminHandler.Students))

	r.mux.HandleFunc("GET /admin/codes", r.middleware.AdminRequire(adminHandler.RegistrationCodes))
	r.mux.HandleFunc("POST /admin/codes/generate", r.middleware.AdminRequire(adminHandler.CreateCode))
	r.mux.HandleFunc("POST /admin/codes/{code}/revoke", r.middleware.AdminRequire(adminHandler.RevokeCode))

	// Маршруты для преподавателя
	r.mux.HandleFunc("GET /teacher/dashboard", r.middleware.TeacherRequire(teacherHandler.Dashboard))
	r.mux.HandleFunc("GET /teacher/disciplines/new", r.middleware.TeacherRequire(teacherHandler.CreateDisciplineForm))
	r.mux.HandleFunc("POST /teacher/disciplines/create", r.middleware.TeacherRequire(teacherHandler.CreateDiscipline))
	r.mux.HandleFunc("GET /teacher/disciplines/{id}/edit", r.middleware.TeacherRequire(teacherHandler.EditDisciplineForm))
	r.mux.HandleFunc("POST /teacher/disciplines/{id}/edit", r.middleware.TeacherRequire(teacherHandler.EditDiscipline))
	r.mux.HandleFunc("POST /teacher/disciplines/{id}/delete", r.middleware.TeacherRequire(teacherHandler.DeleteDiscipline))

	r.mux.HandleFunc("GET /teacher/disciplines/{disciplineId}/labs", r.middleware.TeacherRequire(teacherHandler.LabsList))
	r.mux.HandleFunc("GET /teacher/disciplines/{disciplineId}/labs/new", r.middleware.TeacherRequire(teacherHandler.CreateLabForm))
	r.mux.HandleFunc("POST /teacher/disciplines/{disciplineId}/labs/create", r.middleware.TeacherRequire(teacherHandler.CreateLab))
	r.mux.HandleFunc("GET /teacher/labs/{labId}/edit", r.middleware.TeacherRequire(teacherHandler.EditLabForm))
	r.mux.HandleFunc("POST /teacher/labs/{labId}/edit", r.middleware.TeacherRequire(teacherHandler.EditLab))
	r.mux.HandleFunc("POST /teacher/labs/{labId}/delete", r.middleware.TeacherRequire(teacherHandler.DeleteLab))
}

func (r *Router) Handler() http.Handler {
	handler := r.middleware.Logging(r.handler)
	handler = r.middleware.FlashHandler(handler)
	return handler
}
