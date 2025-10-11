package handlers

import (
	"log"
	"net/http"

	"scavenger/internal/alerts"
	"scavenger/internal/models"
	"scavenger/views"
)

func (h *Handler) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	groups, err := h.db.GetAllGroups()
	if err != nil {
		alerts.FlashWarning(w,r,"Группы не загруженны")
		log.Printf("Failed to get groups: %v", err)
		groups = []models.Group{}
	}

	stats := models.AdminStats{
		TotalReports:   10,
		PendingReports: 3,
		GradedReports:  7,
		TotalGroups:  len(groups),
	}
	reports := []models.LabReport{}

	disciplines, err := h.db.GetDisciplinesWithGroup()
	if err != nil {
		alerts.FlashWarning(w,r,"Дисциплины не загруженны")
		log.Printf("Failed to get disciplines: %v", err)
		disciplines = []models.Discipline{}
	}

	for i, discipline := range disciplines {
		labs,err  := h.db.GetDisciplineLabs(discipline.ID)
		if err != nil {
			continue
		}
		disciplines[i].Labs = append(disciplines[i].Labs, labs...)
	}

	views.AdminDashboard(&stats, reports, disciplines).Render(r.Context(), w)
}
