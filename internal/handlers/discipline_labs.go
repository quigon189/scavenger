package handlers

import (
	"log"
	"net/http"
	"scavenger/internal/alerts"
	"scavenger/views"
	"strconv"
)

func (h *Handler) DisciplineLabs(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusNotFound)
		return
	}

	discipline, err := h.db.GetDisciplineByID(id)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка чтения дисциплины из БД")
		log.Printf("Failed to get disciplint from db: %v", err)
		http.Redirect(w, r, "/admin/disciplines", http.StatusInternalServerError)
		return
	}

	views.DisciplineLabs(discipline).Render(r.Context(), w)
}
