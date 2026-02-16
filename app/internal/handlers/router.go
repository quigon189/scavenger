package handlers

import (
	"net/http"
	"scavenger/internal/middlewares"
	"scavenger/internal/models"
	"scavenger/internal/services"
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
	 return r.Context().Value("session_id").(*models.UserSession)
}

func (r *Router) registerRoutes() {
	authHandler := NewAuthHandler(r.services)
	AdminHandler := NewAdminHandler(r.services)

	r.mux.HandleFunc("/", r.middleware.ActiveRequire(Home))
	r.mux.HandleFunc("/pending", r.middleware.SessionRequire(PendingPage))
	r.mux.HandleFunc("/login", authHandler.Login)
	r.mux.HandleFunc("POST /logout", r.middleware.SessionRequire(authHandler.Logout))
	r.mux.HandleFunc("/register", authHandler.Register)

	r.mux.HandleFunc("/admin/dashboard", r.middleware.AdminRequire(AdminHandler.Dashboard))
	r.mux.HandleFunc("GET /admin/groups", r.middleware.AdminRequire(AdminHandler.Groups))
	r.mux.HandleFunc("GET /admin/teachers", r.middleware.AdminRequire(AdminHandler.Teachers))
	r.mux.HandleFunc("POST /admin/groups/{id}/delete", r.middleware.AdminRequire(AdminHandler.DeleteGroup))
	r.mux.HandleFunc("POST /admin/groups/create", r.middleware.AdminRequire(AdminHandler.CreateGroup))
}

func (r *Router) Handler() http.Handler {
	handler := r.middleware.Logging(r.handler)
	handler = r.middleware.FlashHandler(handler)
	return handler
}
