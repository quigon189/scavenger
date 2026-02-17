package handlers

import (
	"log"
	"net/http"
	"strconv"

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

	BaseWithNavbar(w, r, "Админ-панель", admin.Dashboard(session.User, stats))
}

func (h *AdminHandler) Groups(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)

	groups, err := h.services.Data.GetAllGroups(r.Context())
	if err != nil {
		h.services.Session.FlashError(w, r, "Ошибка при получении групп", err)
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
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

	BaseWithNavbar(w,r,"Управление преподавателями", admin.TeachersPage(teachers))
}

func (h *AdminHandler) CreateTeacher(w http.ResponseWriter, r *http.Request) {

	req := models.RegisterRequest{
		Name: r.FormValue("name"),
		Username: r.FormValue("username"),
		Email: r.FormValue("email"),	
		Role: "teacher",	
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
