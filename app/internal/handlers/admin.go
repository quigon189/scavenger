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

	stats, err := h.services.Data.GetAdminDashboard(r.Context(), session.ID)
	if err != nil {
		log.Printf("WARN Failed to get admin Dashboard stats: %v", err)
	}

	BaseWithNavbar(w, r, "Админ-панель", admin.Dashboard(session.User, stats))
}

func (h *AdminHandler) Groups(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)

	groups, err := h.services.Data.GetGroups(r.Context())
	if err != nil {
		h.session.FlashError(w, r, "Ошибка при получении групп")
		log.Printf("ERR Failed to get groups: %v", err)
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	BaseWithNavbar(w, r, "Управление группами", admin.GroupsPage(session.User, groups))
}

func (h *AdminHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)
	name := r.FormValue("name")

	err := h.services.Data.CreateGroup(r.Context(), session.ID, name)
	if err != nil {
		h.session.FlashError(w, r, err.Error())
		http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
		return
	}

	h.session.FlashSuccess(w, r, "Группа создана")
	http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
}

func (h *AdminHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.session.FlashError(w, r, "Неверный ID группы")
		http.Redirect(w, r, "/admin/gruops", http.StatusSeeOther)
		return
	}

	err = h.services.Data.DeleteGroup(r.Context(), session.ID, id)
	if err != nil {
		h.session.FlashError(w, r, err.Error())
		http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
		return
	}

	h.session.FlashSuccess(w, r, "Группа удалена")
	http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
}

func (h *AdminHandler) Teachers(w http.ResponseWriter, r *http.Request) {
	//_ = getSession(r)

	teachers := []models.User{
		models.User{
			ID: 123,
			Username: "test",
			Name: "Test Test",
			Email: "test@test",
			Status: "active",
		},
	}

	BaseWithNavbar(w,r,"Управление преподавателями", admin.TeachersPage(teachers))
}
