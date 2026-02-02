package handlers

import (
	"net/http"
	"time"
	"web-service/internal/models"
	"web-service/internal/session"
	"web-service/internal/views/pages"
)

type AuthHandler struct {
	session session.SessionService
}

func NewAuthHandler(session session.SessionService) *AuthHandler {
	return &AuthHandler{session: session}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var prevLogin string
	alert := models.Alert{}
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "admin" && password == "password" {
			h.session.SetSession(w, "123", time.Now().Add(3600*time.Second))
			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		} else {
			prevLogin = username
			alert.Type = models.AlertError
			alert.Msg = "Неверно введен логин или пароль"
		}
	}

	pages.Login(prevLogin, alert).Render(r.Context(), w)
}
