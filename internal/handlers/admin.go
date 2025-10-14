package handlers

import (
	"log"
	"net/http"
	"sort"

	"scavenger/internal/alerts"
	"scavenger/internal/models"
	"scavenger/views"
)

func (h *Handler) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	groups, err := h.db.GetAllGroups()
	if err != nil {
		alerts.FlashWarning(w, r, "Группы не загруженны")
		log.Printf("Failed to get groups: %v", err)
		groups = []models.Group{}
	}

	reports, err := h.db.GetAllReports()
	if err != nil {
		alerts.FlashWarning(w, r, "Отчеты не загруженны")
		log.Printf("Failed to get reports: %v", err)
		reports = []models.LabReport{}
	}

	var pendingReports int

	for _, report := range reports {
		if report.Status == "submitted" {
			pendingReports++
		}
	}

	disciplines, err := h.db.GetDisciplinesWithGroup()
	if err != nil {
		alerts.FlashWarning(w, r, "Дисциплины не загруженны")
		log.Printf("Failed to get disciplines: %v", err)
		disciplines = []models.Discipline{}
	}

	for i, discipline := range disciplines {
		labs, err := h.db.GetDisciplineLabs(discipline.ID)
		if err != nil {
			continue
		}
		disciplines[i].Labs = append(disciplines[i].Labs, labs...)
	}

	stats := models.AdminStats{
		TotalReports:     len(reports),
		PendingReports:   pendingReports,
		TotalDisciplines: len(disciplines),
		TotalGroups:      len(groups),
	}

	sort.Slice(reports, func(i, j int) bool {
		return reports[i].UploadedAt.After(reports[j].UploadedAt)
	})

	if len(reports) > 5 {
		reports = reports[:5]
	}

	views.AdminDashboard(&stats, reports, disciplines).Render(r.Context(), w)
}
