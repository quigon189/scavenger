package handlers

import (
	"net/http"
	"scavenger/internal/alerts"
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
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w,r,"/404", http.StatusSeeOther)
	}
	
	h.Dashboard(w,r)
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	role := h.authService.GetUserRole(r)
	switch role {
	case "admin":
		h.AdminDashboard(w,r)
	case "student":
		h.StudentDashboard(w,r)
	default:
		h.Logout(w,r)
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if h.authService.Login(w, r, h.cfg.Users, username, password) {
			http.Redirect(w,r,"/", http.StatusSeeOther)
			return
		} else {
			alerts.FlashError(w,r,"Логин или пароль указанны неверно")
			http.Redirect(w,r,"/login", http.StatusSeeOther)
		}
	}

	views.LoginPage().Render(r.Context(), w)
}

func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	views.NotFound().Render(r.Context(), w)
}

func (h *Handler) Logout(w http.ResponseWriter, r * http.Request) {
	h.authService.Logout(w,r)
	http.Redirect(w,r,"/",http.StatusSeeOther)
}

func (h *Handler) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	views.AdminDashboard().Render(r.Context(), w)
}
