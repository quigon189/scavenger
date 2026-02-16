package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"scavenger/internal/models"
	"scavenger/internal/services"
	"scavenger/internal/views/pages"
)

type AuthHandler struct {
	services *services.Services
}

func NewAuthHandler(services *services.Services) *AuthHandler {
	return &AuthHandler{
		services: services,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	session, _ := h.services.Session.GetSession(r)
	if session != nil && session.ID != "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	var prevLogin string
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		userSession, err := h.services.Auth.Login(r.Context(), &models.LoginRequest{
			Username: username,
			Password: password,
		})
		if err == nil {
			h.services.Session.SetSessionCockie(w, r, userSession.ID, userSession.ExpiresAt)
			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		} else {
			prevLogin = username
			h.services.Session.FlashError(w, r, "Не верный логин или пароль")
			log.Printf("WARN Failed to login user %s: %v", username, err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	Base(w, r, "Вход", pages.Login(prevLogin))
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := h.services.Session.GetSession(r)
	if err == nil {
		h.services.Session.DeleteSession(w, r)
		h.services.Auth.Logout(r.Context(), session.ID)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	session, _ := h.services.Session.GetSession(r)
	if session != nil && session.ID != "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		name := r.FormValue("name")
		email := r.FormValue("email")
		groupID, _ := strconv.Atoi(r.FormValue("group_id"))

		user, err := h.services.Auth.Register(
			r.Context(),
			&models.RegisterRequest{
				Username: username,
				Password: password,
				Name: name,
				Email: email,
				Role: "student",
				GroupID: &groupID,
			},
		)
		if err != nil {
			h.services.Session.FlashError(w, r, fmt.Sprintf("Ошибка при создании пользователя: %s", err.Error()))
			http.Redirect(w, r, "/register", http.StatusSeeOther)
		} else {
			h.services.Session.FlashSuccess(w, r, fmt.Sprintf("Добро пожаловать, %s! Теперь вы можете выполнить вход.", user.Name))
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

	}

	groups, err := h.services.Data.GetAllGroups(r.Context())
	if err != nil {
		h.services.Session.FlashError(w, r, "Ошибка при загрузке групп")
		log.Printf("ERR Failed to get groups: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	Base(w, r, "Регистрация", pages.StudentRegistrationPage(groups))
}
