package handlers

import (
	"net/http"
	"web-service/internal/session"
	"web-service/internal/views/pages"
)

func Home(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value("session").(*session.UserSession)

	if session.User.Role == "admin" {
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	BaseWithNavbar(w, r, "Home", pages.Home(session.User))
}

func PendingPage(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value("session").(*session.UserSession)

	BaseWithNavbar(w, r, "Home", pages.PendingPage(session.User))
}
