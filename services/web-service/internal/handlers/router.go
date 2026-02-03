package handlers

import (
	"net/http"
	"web-service/internal/middlewares"
	"web-service/internal/session"
	"web-service/pkg/apiclient"
)

type Router struct {
	session        *session.SessionService
	authClient     *apiclient.AuthClient
	authMuddleware *middlewares.AuthMiddleware
	handler        http.Handler
	mux            *http.ServeMux
}

func NewRouter(session *session.SessionService, authClient *apiclient.AuthClient) *Router {
	router := &Router{
		mux:        http.NewServeMux(),
		session:    session,
		authClient: authClient,
		authMuddleware: middlewares.NewAuthMiddleware(session),
	}

	router.registerRoutes()

	router.handler = router.mux

	return router
}

func (r *Router) registerRoutes() {
	authHandler := NewAuthHandler(r.session, r.authClient)

	r.mux.HandleFunc("/", r.authMuddleware.SessionRequire(Home))
	r.mux.HandleFunc("/login", authHandler.Login)
	r.mux.HandleFunc("POST /logout", r.authMuddleware.SessionRequire(authHandler.Logout))
	r.mux.HandleFunc("/register", authHandler.Register)
}

func (r *Router) Handler() http.Handler {
	return r.handler
}
