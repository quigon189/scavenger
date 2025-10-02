package handlers

import (
	"net/http"
	"scavenger/internal/alerts"
	"scavenger/internal/models"
	"scavenger/views"
)

func (h *Handler) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	stats := models.AdminStats{
		TotalReports: 10,
		PendingReports: 3,
		GradedReports: 7,
		TotalStudents: 3,
	}
	reports := []models.LabReport{}
	disciplines := []models.Discipline{}
	views.AdminDashboard(&stats, reports, disciplines).Render(r.Context(), w)
}

func (h *Handler) GroupManager(w http.ResponseWriter, r *http.Request) {
	groups, err := h.db.GetAllGroups()
	if err != nil {
		alerts.FlashError(w,r,"Ошибка получения групп")
	}

	disciplines := []models.Discipline{}

	views.GroupsManager(groups, disciplines).Render(r.Context(), w)
}
