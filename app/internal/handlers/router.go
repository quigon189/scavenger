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
	 return r.Context().Value("session").(*models.UserSession)
}

func (r *Router) registerRoutes() {
	authHandler := NewAuthHandler(r.services)
	AdminHandler := NewAdminHandler(r.services)

	r.mux.HandleFunc("/", r.middleware.SessionRequire(Home))
	r.mux.HandleFunc("/login", authHandler.Login)
	r.mux.HandleFunc("POST /logout", r.middleware.SessionRequire(authHandler.Logout))

	r.mux.HandleFunc("GET /register", authHandler.EnterCode)
	r.mux.HandleFunc("POST /register", authHandler.VerifyCode)
	r.mux.HandleFunc("GET /register/complete", authHandler.RegisterComplete)
	r.mux.HandleFunc("POST /register/complete", authHandler.RegisterCompletePost)

	r.mux.HandleFunc("/admin/dashboard", r.middleware.AdminRequire(AdminHandler.Dashboard))

	r.mux.HandleFunc("GET /admin/groups", r.middleware.AdminRequire(AdminHandler.Groups))
	r.mux.HandleFunc("POST /admin/groups/{id}/delete", r.middleware.AdminRequire(AdminHandler.DeleteGroup))
	r.mux.HandleFunc("POST /admin/groups/create", r.middleware.AdminRequire(AdminHandler.CreateGroup))

	r.mux.HandleFunc("GET /admin/teachers", r.middleware.AdminRequire(AdminHandler.Teachers))
	r.mux.HandleFunc("POST /admin/teachers/{id}/delete", r.middleware.AdminRequire(AdminHandler.DeleteTeacher))
	r.mux.HandleFunc("POST /admin/teachers/create", r.middleware.AdminRequire(AdminHandler.CreateTeacher))

	r.mux.HandleFunc("GET /admin/students", r.middleware.AdminRequire(AdminHandler.Students))

	r.mux.HandleFunc("GET /admin/codes", r.middleware.AdminRequire(AdminHandler.RegistrationCodes))
	r.mux.HandleFunc("POST /admin/codes/generate", r.middleware.AdminRequire(AdminHandler.CreateCode))
	r.mux.HandleFunc("POST /admin/codes/{code}/revoke", r.middleware.AdminRequire(AdminHandler.RevokeCode))
}

func (r *Router) Handler() http.Handler {
	handler := r.middleware.Logging(r.handler)
	handler = r.middleware.FlashHandler(handler)
	return handler
}
