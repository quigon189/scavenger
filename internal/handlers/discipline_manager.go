package handlers

import (
	"log"
	"net/http"
	"scavenger/internal/alerts"
	"scavenger/internal/models"
	"scavenger/views"
	"strconv"
)

func (h *Handler) DisciplinesManager(w http.ResponseWriter, r *http.Request) {
	groups, err := h.db.GetAllGroupsWithDisciplines()
	if err != nil {
		alerts.FlashError(w, r, "Ошибка при получении групп")
		log.Printf("Failed to get groups: %v", err)
		groups = []models.Group{}
	}

	disciplines, err := h.db.GetDisciplines()
	if err != nil {
		alerts.FlashError(w, r, "Ошибка при получении дисциплин")
		log.Printf("Failed to get disciplines: %v", err)
		disciplines = []models.Discipline{}
	}

	views.DisciplinesManager(disciplines, groups).Render(r.Context(), w)
}

func (h *Handler) AddDiscipline(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		alerts.FlashError(w, r, "Ошибка обработки формы")
		log.Printf("Failed to parse form: %v", err)
		http.Redirect(w, r, "/admin/disciplines", http.StatusSeeOther)
		return
	}

	discName := r.FormValue("name")
	groupID := r.FormValue("group")

	disc := &models.Discipline{
		Name: discName,
	}

	if groupID != "" {
		gID, err := strconv.Atoi(groupID)
		if err != nil {
			alerts.FlashError(w, r, "Ошибка чтения ID группы")
			log.Printf("Failed to conv group id: %v", err)
			http.Redirect(w, r, "/admin/disciplines", http.StatusSeeOther)
			return
		}
		disc.GroupID = &gID
	}

	err = h.db.CreateDiscipline(disc)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка создания дисциплины")
		log.Printf("Failed to create discipline: %v", err)
		http.Redirect(w, r, "/admin/disciplines", http.StatusSeeOther)
		return
	}

	alerts.FlashSuccess(w, r, "Дисциплина создана")
	http.Redirect(w, r, "/admin/disciplines", http.StatusSeeOther)
}
