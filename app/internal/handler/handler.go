package handler

import (
	"log/slog"
	"net/http"

	"scavenger/internal/config"
	"scavenger/internal/csrf"
	"scavenger/internal/domain"
	"scavenger/internal/routes"
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
	r.Use(Recoverer)
	r.Use(csrf.Middleware)
	r.Use(h.AuthMiddleware)

	fs := http.FileServer(http.Dir("web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Group(func(pub chi.Router) {
		pub.Get(routes.Login, h.LoginPage)
		pub.Post(routes.Login, h.LoginSubmit)
		pub.Post(routes.Logout, h.Logout)
	})

	r.Group(func(g chi.Router) {
		g.Use(csrf.Verify)
		g.Use(RequireAuth)
		g.Get(routes.Home, h.Home)
		g.Get(routes.DisciplineID.Pattern, h.Discipline)
		g.Group(func(r chi.Router) {
			r.Use(RequireRole(domain.RoleTeacher))
			r.Get(routes.MaterialNew.Pattern, h.MaterialFormNew)
			r.Post(routes.MaterialNew.Pattern, h.MaterialFormNewSubmit)

			//r.Get(routes.MaterialEdit.Pattern, h.MaterialFormEdit)
			//r.Post(routes.MaterialEdit.Pattern, h.MaterialFormEditSubmit)

			r.Post(routes.MaterialAttachment, h.MaterialAttachmentUpload)
			r.Post(routes.MaterialAttachmentDelete.Pattern, h.MaterialAttachmentDelete)


		})
		g.Get(routes.MaterialID.Pattern, h.MaterialView)
		g.Post(routes.MaterialMove.Pattern, h.MaterialMove)
	})

	return r
}
