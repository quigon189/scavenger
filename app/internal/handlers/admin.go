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
		h.services.Session.FlashError(w, r, "Ошибка при получении групп")
		log.Printf("ERR Failed to get groups: %v", err)
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	BaseWithNavbar(w, r, "Управление группами", admin.GroupsPage(session.User, groups))
}

func (h *AdminHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")

	err := h.services.Data.CreateGroup(r.Context(), &models.Group{Name: name})
	if err != nil {
		h.services.Session.FlashError(w, r, err.Error())
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
		h.services.Session.FlashError(w, r, "Неверный ID группы")
		http.Redirect(w, r, "/admin/gruops", http.StatusSeeOther)
		return
	}

	err = h.services.Data.DeleteGroup(r.Context(), id)
	if err != nil {
		h.services.Session.FlashError(w, r, err.Error())
		http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
		return
	}

	h.services.Session.FlashSuccess(w, r, "Группа удалена")
	http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
}

func (h *AdminHandler) Teachers(w http.ResponseWriter, r *http.Request) {

	teachers, err := h.services.Data.GetAllTeachers(r.Context())
	if err != nil {
		h.services.Session.FlashError(w, r, err.Error())
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	BaseWithNavbar(w,r,"Управление преподавателями", admin.TeachersPage(teachers))
}
