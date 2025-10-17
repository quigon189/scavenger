package handlers

import (
	"log"
	"net/http"
	"scavenger/internal/alerts"
	"scavenger/internal/auth"
	"scavenger/internal/database"
	filestorage "scavenger/internal/file_storage"
	"scavenger/internal/models"
	"scavenger/views"
)

type Handler struct {
	authService *auth.AuthService
	cfg         *models.Config
	db          *database.Database
	fs          *filestorage.FileStorage
}

func NewHandler(cfg *models.Config, db *database.Database, fs *filestorage.FileStorage) *Handler {
	return &Handler{
		authService: auth.New(cfg.Auth),
		cfg:         cfg,
		db:          db,
		fs:          fs,
	}
}
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
	}

	h.Dashboard(w, r)
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	role := h.authService.GetUserRole(r)
	switch role {
	case "admin":
		h.AdminDashboard(w, r)
	case "student":
		h.StudentDashboard(w, r)
	default:
		h.Logout(w, r)
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err := h.db.GetUserByUsername(username)
		if err != nil {
			alerts.FlashError(w, r, "Логин или пароль указанны неверно")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if user.RoleName == string(models.StudentRole) {
			err := h.db.GetStudentGroup(user)
			if err != nil {
				alerts.FlashError(w, r, "Ошибка при получении группы пользователя")
				log.Printf("Failed to det student group: %v", err)
				return
			}
		}

		if h.authService.Login(w, r, user, password) {
			log.Printf("Success login user :%v", user)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else {
			alerts.FlashError(w, r, "Логин или пароль указанны неверно")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	views.LoginPage().Render(r.Context(), w)
}

func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	views.NotFound().Render(r.Context(), w)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	h.authService.Logout(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	user := models.GetUserFromContext(r.Context())

	if err := r.ParseForm(); err != nil {
		alerts.FlashError(w, r, "Ошибка обработки формы")
		log.Printf("Failed to parse form: %v", err)
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	currentPassword := r.FormValue("current_password")
	newPassword := r.FormValue("new_password")
	confirmPassword := r.FormValue("confirm_password")

	if currentPassword == "" || newPassword == "" || confirmPassword == "" {
		alerts.FlashError(w, r, "Заполнены не все поля")
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	if newPassword != confirmPassword || len(newPassword) < 6 {
		alerts.FlashError(w, r, "Указан некорректный пароль")
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	user, err := h.db.GetUserByUsername(user.Username)
	if err != nil {
		alerts.FlashError(w, r, "Пользователь не найден")
		log.Printf("Failed to get user: %v", err)
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	if err = auth.ComparePassword(*user, currentPassword); err != nil {
		alerts.FlashError(w, r, "Пароль указан неверно")
		log.Printf("Failed to compare: %v", err)
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	passHash, err := auth.GeneratePassHash(newPassword)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка обработки пароля")
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	user.PasswordHash = passHash

	if err := h.db.UpdateUser(user); err != nil {
		alerts.FlashError(w, r, "Ошибка обновления данных пользователя")
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	alerts.FlashSuccess(w, r, "Пароль обнавлен")
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}
