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
	dataClient *apiclient.DataClient
}

func NewAuthHandler(session *session.SessionService, authClient *apiclient.AuthClient, dataClient *apiclient.DataClient) *AuthHandler {
	return &AuthHandler{
		session:    session,
		authClient: authClient,
		dataClient: dataClient,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	_, sessionID, _ := h.session.GetSession(r)
	if sessionID != "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	var prevLogin string
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
			h.session.SetSession(w, r, resp.SessionID, resp.ExpiresAt)
			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		} else {
			prevLogin = username
			h.session.FlashError(w, r, "Не верный логин или пароль")
			log.Printf("WARN Failed to login user %s: %v", username, err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	Base(w, r, "Вход", pages.Login(prevLogin))
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	_, sessionID, err := h.session.GetSession(r)
	if err == nil {
		h.session.DeleteSession(w, r)
		h.authClient.Logout(r.Context(), sessionID)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.session.FlashSuccess(w, r, "Вы зарегистрированны!")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	groups := []models.Group{}
	gs, err := h.dataClient.GetGroups(r.Context())
	if err == nil {
		for _, group := range gs {
			groups = append(groups, models.Group{
				ID:   group.ID,
				Name: group.Name,
			})
		}
	}

	Base(w, r, "Регистрация", pages.StudentRegistrationPage(groups))
}
