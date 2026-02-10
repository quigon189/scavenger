package handlers

import (
	"log"
	"net/http"
	"strconv"

	"web-service/internal/services"
	"web-service/internal/session"
	"web-service/internal/views/pages/admin"
)

type AdminHandler struct {
	session  *session.SessionService
	services *services.Services
}

func NewAdminHandler(session *session.SessionService, services *services.Services) *AdminHandler {
	return &AdminHandler{
		session:  session,
		services: services,
	}
}

func (h *AdminHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value("session").(*session.UserSession)

	stats, err := h.services.Data.GetAdminDashboard(r.Context(), session.ID)
	if err != nil {
		log.Printf("WARN Failed to get admin Dashboard stats: %v", err)
	}

	BaseWithNavbar(w, r, "Админ-панель", admin.Dashboard(session.User, stats))
}

func (h *AdminHandler) Groups(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*session.UserSession)

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
	session := r.Context().Value("session_id").(*session.UserSession)
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
	session := r.Context().Value("session_id").(*session.UserSession)
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
