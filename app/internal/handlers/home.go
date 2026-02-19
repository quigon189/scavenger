package handlers

import (
	"net/http"
	"scavenger/internal/models"
	"scavenger/internal/views/pages"
)

func Home(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value("session").(*models.UserSession)

	switch session.User.Role {
	case "admin":
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
	case "teacher":
		http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
	default:
		BaseWithNavbar(w, r, "Home", pages.Home(session.User))
	}

	BaseWithNavbar(w, r, "Home", pages.Home(session.User))
}
