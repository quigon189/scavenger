package handlers

import (
	"log"
	"net/http"
	"web-service/internal/models"
	"web-service/internal/session"
	"web-service/internal/views/pages"
	"web-service/pkg/apiclient"
)

type AuthHandler struct {
	session    *session.SessionService
	authClient *apiclient.AuthClient
}

func NewAuthHandler(session *session.SessionService, authClient *apiclient.AuthClient) *AuthHandler {
	return &AuthHandler{session: session, authClient: authClient}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var prevLogin string
	alert := models.Alert{}
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		resp, err := h.authClient.Login(
			r.Context(),
			apiclient.LoginRequest{
				Username: username,
				Password: password,
			},
		)
		if err == nil && resp.SessionID != "" {
			h.session.SetSession(w, resp.SessionID, resp.ExpiresAt)
			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		} else {
			prevLogin = username
			alert.Type = models.AlertError
			alert.Msg = "Неверно введен логин или пароль"
			log.Printf("WARN Failed to login user %s: %v", username, err)
		}
	}

	pages.Login(prevLogin, alert).Render(r.Context(), w)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	_, sessionID, err := h.session.GetSession(r)
	if err == nil {
		h.session.DeleteSession(w)
		h.authClient.Logout(r.Context(), sessionID)
	}
	http.Redirect(w,r,"/", http.StatusSeeOther)
}
