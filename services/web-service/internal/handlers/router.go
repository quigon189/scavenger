package handlers

import (
	"net/http"
	"web-service/internal/middlewares"
	"web-service/internal/session"
	"web-service/pkg/apiclient"
)

type Router struct {
	session    *session.SessionService
	authClient *apiclient.AuthClient
	dataClient *apiclient.DataClient
	middleware *middlewares.Middleware
	handler    http.Handler
	mux        *http.ServeMux
}

func NewRouter(session *session.SessionService, authClient *apiclient.AuthClient, dataClient *apiclient.DataClient) *Router {
	router := &Router{
		mux:        http.NewServeMux(),
		session:    session,
		authClient: authClient,
		dataClient: dataClient,
		middleware: middlewares.NewMiddleware(session),
	}

	router.registerRoutes()

	router.handler = router.mux

	return router
}

func (r *Router) registerRoutes() {
	authHandler := NewAuthHandler(r.session, r.authClient, r.dataClient)

	r.mux.HandleFunc("/", r.middleware.ActiveRequire(Home))
	r.mux.HandleFunc("/pending", r.middleware.SessionRequire(PendingPage))
	r.mux.HandleFunc("/login", authHandler.Login)
	r.mux.HandleFunc("POST /logout", r.middleware.SessionRequire(authHandler.Logout))
	r.mux.HandleFunc("/register", authHandler.Register)
}

func (r *Router) Handler() http.Handler {
	handler := r.middleware.Logging(r.handler)
	handler = r.middleware.FlashHandler(handler)
	return handler
}
