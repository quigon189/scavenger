package handlers

import (
	"net/http"
	"scavenger/internal/auth"
	"scavenger/internal/models"
	"scavenger/views"
)

type Handler struct {
	authService *auth.AuthService
	cfg *models.Config
}

func NewHandler(cfg *models.Config) *Handler {
	return &Handler{
		authService: auth.New(cfg.Auth),
		cfg: cfg,
	}
}

func (h *Handler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !h.authService.IsAuthenticated(r) {
			http.Redirect(w,r,"/login", http.StatusSeeOther)
		}

		next.ServeHTTP(w,r)
	}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	role := h.authService.GetUserRole(r)

	switch role {
	case "admin":
		http.Redirect(w, r , "/admin/dashboard", http.StatusSeeOther)
	case "student":
		http.Redirect(w, r, "/student/dashboard", http.StatusSeeOther)
	default:
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if h.authService.Login(w, r, h.cfg.Users, username, password) {
			http.Redirect(w,r,"/", http.StatusSeeOther)
			return
		}


	}

	views.LoginPage().Render(r.Context(), w)
}
