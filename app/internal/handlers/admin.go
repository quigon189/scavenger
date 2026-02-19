package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"scavenger/internal/models"
	"scavenger/internal/services"
	"scavenger/internal/views/pages/admin"
)

type AdminHandler struct {
	services *services.Services
}

func NewAdminHandler(services *services.Services) *AdminHandler {
	return &AdminHandler{
		services: services,
	}
}

func (h *AdminHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)

	stats, err := h.services.Data.GetAdminDashboard(r.Context())
	if err != nil {
		log.Printf("WARN Failed to get admin Dashboard stats: %v", err)
	}

	codes, _ := h.services.Auth.GetAllCodes(r.Context())
	if len(codes) > 5 {
		codes = codes[:5]
	}

	BaseWithNavbar(w, r, "Админ-панель", admin.Dashboard(session.User, stats, codes))
}

func (h *AdminHandler) Groups(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)

	groups, err := h.services.Data.GetAllGroups(r.Context())
	if err != nil {
		h.services.Session.FlashError(w, r, "Ошибка при получении групп", err)
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	for i := range groups {
		students, _ := h.services.Data.GetStudentsByGroup(r.Context(), groups[i].ID)
		groups[i].Students = append(groups[i].Students, students...)
	}

	BaseWithNavbar(w, r, "Управление группами", admin.GroupsPage(session.User, groups))
}

func (h *AdminHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")

	err := h.services.Data.CreateGroup(r.Context(), &models.Group{Name: name})
	if err != nil {
		h.services.Session.FlashError(w, r, "Ошибка при создании группы", err)
		http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
		return
	}

	h.services.Session.FlashSuccess(w, r, "Группа создана")
	http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
}

func (h *AdminHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.services.Session.FlashError(w, r, "Неверный ID группы", err)
		http.Redirect(w, r, "/admin/gruops", http.StatusSeeOther)
		return
	}

	err = h.services.Data.DeleteGroup(r.Context(), id)
	if err != nil {
		h.services.Session.FlashError(w, r, "Ошибка удаления группы", err)
		http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
		return
	}

	h.services.Session.FlashSuccess(w, r, "Группа удалена")
	http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
}

func (h *AdminHandler) Teachers(w http.ResponseWriter, r *http.Request) {
	teachers, err := h.services.Data.GetAllTeachers(r.Context())
	if err != nil {
		h.services.Session.FlashError(w, r, "Ошибка при получении пользователей", err)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	BaseWithNavbar(w, r, "Управление преподавателями", admin.TeachersPage(teachers))
}

func (h *AdminHandler) CreateTeacher(w http.ResponseWriter, r *http.Request) {

	req := models.RegisterRequest{
		Name:     r.FormValue("name"),
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Role:     "teacher",
		Password: r.FormValue("password"),
	}

	user, err := h.services.Auth.Register(r.Context(), &req)
	if err != nil {
		h.services.Session.FlashError(w, r, "Ошибка регистрации пользователя", err)
		http.Redirect(w, r, "/admin/teachers", http.StatusSeeOther)
		return
	}
	err = h.services.Data.CreateTeacher(r.Context(), &models.Teacher{ID: user.ID})
	if err != nil {
		h.services.Session.FlashError(w, r, "Ошибка создания преподавателя", err)
		http.Redirect(w, r, "/admin/teachers", http.StatusSeeOther)
		return
	}

	h.services.Session.FlashSuccess(w, r, "Учетная запись преподавателя создана")
	http.Redirect(w, r, "/admin/teachers", http.StatusSeeOther)
}

func (h *AdminHandler) DeleteTeacher(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.services.Session.FlashError(w, r, "Неверный ID преподавателя", err)
		http.Redirect(w, r, "/admin/teachers", http.StatusSeeOther)
		return
	}

	err = h.services.Data.DeleteTeacher(r.Context(), id)
	if err != nil {
		h.services.Session.FlashError(w, r, "Ошибка удаления преподавателя", err)
		http.Redirect(w, r, "/admin/teachers", http.StatusSeeOther)
		return
	}

	h.services.Session.FlashSuccess(w, r, "Группа удалена")
	http.Redirect(w, r, "/admin/teachers", http.StatusSeeOther)

}

func (h *AdminHandler) Students(w http.ResponseWriter, r *http.Request) {
	students, err := h.services.Data.GetAllStudents(r.Context())
	if err != nil {
		log.Printf("WARN Failed to get students: %v", err)
		students = []models.Student{}
	}

	BaseWithNavbar(w, r, "Управление студентами", admin.StudentsPage(students))
}

func (h *AdminHandler) RegistrationCodes(w http.ResponseWriter, r *http.Request) {
	codes, err := h.services.Auth.GetAllCodes(r.Context())
	if err != nil {
		log.Printf("WARN Failed to get registration codes: %v", err)
		codes = []models.RegistrationCode{}
	}

	groups, _ := h.services.Data.GetAllGroups(r.Context())

	BaseWithNavbar(w, r, "Коды регистрации", admin.RegistrationCodesPage(codes, groups))
}

func (h *AdminHandler) CreateCode(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	email := r.FormValue("email")
	role := r.FormValue("role")
	groupIDStr := r.FormValue("group_id")

	var groupID *int
	if role == "student" {
		if groupIDStr != "" {
			id, err := strconv.Atoi(groupIDStr)
			if err == nil {
				groupID = &id
			}
		} else {
			h.services.Session.FlashError(w, r, "не выбрана группа", fmt.Errorf("group not set"))
			http.Redirect(w, r, "/admin/codes", http.StatusSeeOther)
			return
		}
	}

	code := &models.RegistrationCode{
		Name:      name,
		Email:     email,
		Role:      role,
		GroupID:   groupID,
		ExpiresAt: time.Now().Add(24 * 7 * time.Hour),
	}

	err := h.services.Auth.GenerateCode(r.Context(), code)
	if err != nil {
		log.Printf("ERROR Failed to generate code: %v", err)
		h.services.Session.FlashError(w, r, "Ошибка при генерации кода: ", err)
	} else {
		// Сохраняем сгенерированный код в сессии, чтобы показать в модалке
		h.services.Session.FlashSuccess(w, r, "Код успешно сгенерирован: " + code.Code)
	}

	http.Redirect(w, r, "/admin/codes", http.StatusSeeOther)
}

func (h *AdminHandler) RevokeCode(w http.ResponseWriter, r *http.Request) {

    code := r.PathValue("code")
    if code == "" {
        h.services.Session.FlashError(w, r, "Код не указан", fmt.Errorf("bad code"))
        http.Redirect(w, r, "/admin/codes", http.StatusSeeOther)
        return
    }

	err := h.services.Auth.RevokeCode(r.Context(), code)
    if err != nil {
        log.Printf("ERROR Failed to revoke code %s: %v", code, err)
        h.services.Session.FlashError(w, r, "Ошибка при отзыве кода: "+err.Error(), err)
    } else {
        h.services.Session.FlashSuccess(w, r, "Код успешно отозван")
    }

    http.Redirect(w, r, "/admin/codes", http.StatusSeeOther)
}
