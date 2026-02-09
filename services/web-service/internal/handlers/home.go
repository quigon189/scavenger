package handlers

import (
	"net/http"
	"web-service/internal/session"
	"web-service/internal/views/pages"
)

func Home(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value("session").(*session.UserSession)
	if !ok {
		http.Error(w,"Failed to get session", http.StatusInternalServerError)
		return
	}

	BaseWithNavbar(w,r,"Home", pages.Home(session.User))	
}

func PendingPage(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value("session").(*session.UserSession)

	BaseWithNavbar(w,r,"Home", pages.PendingPage(session.User))
}
