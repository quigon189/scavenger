package handler

import (
	"log/slog"
	"net/http"

	"scavenger/internal/config"
	"scavenger/internal/csrf"
	"scavenger/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Deps struct {
	Config   *config.Config
	Services *service.Services
}

type Handler struct {
	cfg *config.Config
	svc *service.Services
}

func New(d Deps) *Handler {
	return &Handler{
		cfg: d.Config,
		svc: d.Services,
	}
}

func (h *Handler) Router() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(RequestLogger(slog.Default()))
	r.Use(middleware.Recoverer)
	r.Use(csrf.Middleware)

	fs := http.FileServer(http.Dir("web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Group(func(pub chi.Router) {
		pub.Get("/login", h.LoginPage)
		pub.Post("/login", h.LoginSubmit)
		pub.Post("/logout", h.Logout)
	})

	r.Group(func(g chi.Router) {
		g.Use(csrf.Middleware)
		g.Use(h.AuthMiddleware)
		g.Get("/", h.Home)
	})

	return r
}
