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
			h.services.Session.FlashError(w, r, "Не верный логин или пароль", err)
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

func (h *AuthHandler) EnterCode(w http.ResponseWriter, r *http.Request) {
	session, _ := h.services.Session.GetSession(r)
	if session != nil && session.ID != "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	Base(w, r, "Ввод кода", pages.EnterCode())
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
				Name:     name,
				Email:    email,
				Role:     "student",
				GroupID:  &groupID,
			},
		)
		if err != nil {
			h.services.Session.FlashError(w, r, fmt.Sprintf("Ошибка при создании пользователя: %s", err.Error()), err)
			http.Redirect(w, r, "/register", http.StatusSeeOther)
		} else {
			h.services.Session.FlashSuccess(w, r, fmt.Sprintf("Добро пожаловать, %s! Теперь вы можете выполнить вход.", user.Name))
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

	}

	groups, err := h.services.Data.GetAllGroups(r.Context())
	if err != nil {
		h.services.Session.FlashError(w, r, "Ошибка при загрузке групп", err)
		log.Printf("ERR Failed to get groups: %v", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	Base(w, r, "Регистрация", pages.StudentRegistrationPage(groups))
}

func (h *AuthHandler) VerifyCode(w http.ResponseWriter, r *http.Request) {
	sess, _ := h.services.Session.GetSession(r)
	if sess != nil && sess.ID != "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	code := models.RegCodeRequest{}

	code.Code = r.FormValue("code")
	code.Email = r.FormValue("email")

	regCode, err := h.services.Auth.GetCodeInfo(r.Context(), &code)
	if err != nil {
		h.services.Session.FlashError(w, r, "Неверный код или email", err)
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	h.services.Session.SetCodeCookie(w, r, regCode)

	http.Redirect(w, r, "/register/complete", http.StatusSeeOther)
}

func (h *AuthHandler) RegisterComplete(w http.ResponseWriter, r *http.Request) {
	sess, _ := h.services.Session.GetSession(r)
	if sess != nil && sess.ID != "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	code, err := h.services.Session.GetCodeCookie(w, r)
	if err != nil {
		h.services.Session.FlashError(w, r, "Сначала надо указать код и email", fmt.Errorf("failed to get code cookie"))
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	var groupName string
	if code.GroupID != nil {
		// Получаем название группы по ID
		group, err := h.services.Data.GetGroup(r.Context(), *code.GroupID)
		if err == nil {
			groupName = group.Name
		} else {
			groupName = fmt.Sprintf("Группа %d", *code.GroupID)
		}
	}

	Base(w, r, "Завершение регистрации", pages.RegisterComplete(code.Name, code.Email, groupName, code.Role))
}

func (h *AuthHandler) RegisterCompletePost(w http.ResponseWriter, r *http.Request) {
	sess, _ := h.services.Session.GetSession(r)
	if sess != nil && sess.ID != "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	code, err := h.services.Session.GetCodeCookie(w, r)
	if err != nil {
		h.services.Session.FlashError(w, r, "Сначала надо указать код и email", fmt.Errorf("failed to get code cookie"))
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	req := &models.RegisterRequest{
		Username: username,
		Password: password,
		Name:     code.Name,
		Email:    code.Email,
		Role:     code.Role,
		GroupID:  code.GroupID,
	}

	user, err := h.services.Auth.Register(r.Context(), req)
	if err != nil {
		h.services.Session.FlashError(w, r, fmt.Sprintf("Ошибка регистрации: %s", err.Error()), err)
		http.Redirect(w, r, "/register/complete", http.StatusSeeOther)
		return
	}

	if user.Role == "student" {
		h.services.Data.CreateStudent(r.Context(), &models.Student{
			ID:      user.ID,
			GroupID: *code.GroupID,
		})
	} else if user.Role == "teacher" {
		h.services.Data.CreateTeacher(r.Context(), &models.Teacher{ID: user.ID})
	}

	h.services.Session.DeleteCodeCookie(w, r)
	h.services.Auth.RevokeCode(r.Context(), code.Code)

	h.services.Session.FlashSuccess(w, r, fmt.Sprintf("Добро пожаловать, %s! Теперь вы можете войти.", user.Name))
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
