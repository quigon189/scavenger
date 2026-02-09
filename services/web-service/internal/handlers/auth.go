package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"web-service/internal/services"
	"web-service/internal/session"
	"web-service/internal/views/pages"
	"web-service/pkg/apiclient"
)

type AuthHandler struct {
	session    *session.SessionService
	authClient *apiclient.AuthClient
	services   *services.Services
}

func NewAuthHandler(session *session.SessionService, authClient *apiclient.AuthClient, dataClient *apiclient.DataClient) *AuthHandler {
	services := services.NewServices(dataClient)
	return &AuthHandler{
		session:    session,
		authClient: authClient,
		services:   services,
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
	_, sessionID, _ := h.session.GetSession(r)
	if sessionID != "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		name := r.FormValue("name")
		groupID, _ := strconv.Atoi(r.FormValue("group_id"))

		req := apiclient.RegisterRequest{
			Username: username,
			Password: password,
			Name:     name,
			GroupID:  &groupID,
			Role:     "student",
		}

		user, err := h.authClient.Register(r.Context(), req)
		if err != nil {
			h.session.FlashError(w, r, "Ошибка при создании пользователя")
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		h.session.FlashSuccess(w, r, fmt.Sprintf("Добро пожаловать, %s!\nТеперь вы можете выполнить вход.", user.Name))
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	groups, err := h.services.Data.GetGroups(r.Context())
	if err != nil {
		h.session.FlashError(w, r, "Ошибка при загрузке групп")
		log.Printf("ERR Failed to get groups: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	Base(w, r, "Регистрация", pages.StudentRegistrationPage(groups))
}
