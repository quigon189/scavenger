package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"data-service/internal/middlewares"
	"data-service/internal/services"
	"data-service/internal/session"
)

type Router struct {
	mux            *http.ServeMux
	handler        http.Handler
	services       *services.Services
	authMiddleware *middlewares.AuthMiddleware
}

func NewRouter(services *services.Services, session *session.SessionService) *Router {
	router := &Router{
		mux:            http.NewServeMux(),
		services:       services,
		authMiddleware: middlewares.NewAuthMiddleware(session),
	}

	router.registerHandlers()

	router.handler = router.wrapMiddleware(router.mux)

	return router
}

func (r *Router) Handler() http.Handler {
	return r.handler
}

func (r *Router) registerHandlers() {

	studentHandlers := NewStudentHandlers(r.services.Students, r.authMiddleware)
	teacherHandlers := NewTeacherHandlers(r.services.Teachers, r.authMiddleware)
	disciplineHandlers := NewDisciplineHandlers(r.services.Disciplines, r.authMiddleware)
	fileHandlers := NewFileHandlers(r.services.Files, r.authMiddleware)
	labHandlers := NewLabHandlers(r.services.Labs, r.authMiddleware)
	reportHandlers := NewReportHandlers(r.services.Reports, r.authMiddleware)
	groupHandlers := NewGroupHandlers(r.services.Groups, r.authMiddleware)
	periodHandlers := NewPeriodHandlers(r.services.Periods, r.authMiddleware)

	studentHandlers.RegisterStudentRoutes(r.mux)
	teacherHandlers.RegisterTeacherRoutes(r.mux)
	disciplineHandlers.RegisterDisciplineRoutes(r.mux)
	fileHandlers.RegisterFileRoutes(r.mux)
	labHandlers.RegisterLabRoutes(r.mux)
	reportHandlers.RegisterReportRoutes(r.mux)
	groupHandlers.RegisterGroupRoutes(r.mux)
	periodHandlers.RegisterPeriodRoutes(r.mux)

	r.mux.HandleFunc("GET /health", r.healthHandler)

	r.mux.HandleFunc("/", r.notFoundHandler)
}

func (r *Router) wrapMiddleware(handler http.Handler) http.Handler {
	return corsMiddleware(
		LoggingMiddleware(
			r.recoveryMiddleware(
				handler,
			),
		),
	)
}

func (r *Router) healthHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":    "healthy",
		"service":   "data-service",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(response)
}

func (r *Router) notFoundHandler(w http.ResponseWriter, req *http.Request) {
	writeError(w, http.StatusNotFound, fmt.Errorf("route not found: %s %s", req.Method, req.URL.Path))
}

func (r *Router) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				writeError(w, http.StatusInternalServerError, fmt.Errorf("internal server error"))
			}
		}()

		next.ServeHTTP(w, req)
	})
}
