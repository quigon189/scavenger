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
	_, sessionID, _ := h.session.GetSession(r)
	if sessionID != "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	var prevLogin string
	alerts := []models.Alert{}
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
			alert := models.Alert{
				Type: models.AlertError,
				Msg:  "Неверно введен логин или пароль",
			}

			alerts = append(alerts, alert)
			log.Printf("WARN Failed to login user %s: %v", username, err)
		}
	}

	Base(w, r, "Login", pages.Login(prevLogin))
	SetAlerts(w, alerts...)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	_, sessionID, err := h.session.GetSession(r)
	if err == nil {
		h.session.DeleteSession(w)
		h.authClient.Logout(r.Context(), sessionID)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	alerts := []models.Alert{}
	if r.Method == http.MethodPost {
		alerts = append(alerts, models.Alert{
			Type: models.AlertSuccess,
			Msg: "You registred!",
		})	

		SetAlerts(w, alerts...)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	groups := []models.Group{}
	groups = append(groups, models.Group{ID: 1, Name: "Test1"})
	groups = append(groups, models.Group{ID: 2, Name: "Test2"})

	Base(w, r, "Регистрация", pages.StudentRegistrationPage(groups))
}
