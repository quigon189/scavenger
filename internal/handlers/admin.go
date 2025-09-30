package handlers

import (
	"net/http"
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
	reports := h.cfg.LabReports
	disciplines := h.cfg.Disciplines
	views.AdminDashboard(&stats, reports, disciplines).Render(r.Context(), w)
}
