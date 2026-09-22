package handler

import (
	"log/slog"
	"net/http"

	"scavenger/internal/config"
	"scavenger/internal/reqlog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Deps struct {
	Config *config.Config
}

type Handler struct {
	cfg *config.Config
}

func New(d Deps) *Handler {
	return &Handler{cfg: d.Config}
}

func (h *Handler) Router() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(RequestLogger(slog.Default()))
	r.Use(middleware.Recoverer)

	fs := http.FileServer(http.Dir("web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Group(func(pub chi.Router) {
		pub.Get("/", func(w http.ResponseWriter, r *http.Request) {
			log := reqlog.From(r.Context())
			log.Info("Hello from slog!")
			w.Write([]byte("Hello!"))
		})
	})

	return r
}
