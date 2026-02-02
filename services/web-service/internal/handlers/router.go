package handlers

import (
	"net/http"
	"web-service/internal/session"
)

type Router struct {
	session *session.SessionService
	handler http.Handler
	mux     *http.ServeMux
}

func NewRouter(session *session.SessionService) *Router {
	router := &Router{
		mux: http.NewServeMux(),
		session: session,
	}

	router.registerRoutes()

	router.handler = router.mux

	return router
}

func (r *Router) registerRoutes() {
	authHandler := NewAuthHandler(*r.session)

	r.mux.HandleFunc("/home", Home)
	r.mux.HandleFunc("/login", authHandler.Login)
}

func (r *Router) Handler() http.Handler {
	return r.handler
}
