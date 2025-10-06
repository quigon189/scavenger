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

	disciplines, err := h.db.GetDisciplinesWithoutGroup()
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

func (h *Handler) EditDiscipline(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusNotFound)
		return
	}

	err = r.ParseForm()
	if err != nil {
		alerts.FlashError(w, r, "Ошибка обработки формы")
		log.Printf("Failed to parse form: %v", err)
		http.Redirect(w, r, "/admin/disciplines", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	var groupID *int
	if r.FormValue("group") != "" {
		gID, err := strconv.Atoi(r.FormValue("group"))
		groupID = &gID
		if err != nil {
			alerts.FlashError(w, r, "Ошибка чтения ID группы")
			log.Printf("Failed to conv group id: %v", err)
			http.Redirect(w, r, "/admin/disciplines", http.StatusSeeOther)
			return
		}

	}

	editDisc, err := h.db.GetDisciplineByID(id)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка чтения дисциплины из БД")
		log.Printf("Failed to get disc from db: %v", err)
		http.Redirect(w, r, "/admin/disciplines", http.StatusSeeOther)
		return
	}

	editDisc.Name = name
	editDisc.GroupID = groupID

	err = h.db.UpdateDiscipline(editDisc)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка обновления дисциплины")
		log.Printf("Failed to update disc: %v", err)
		http.Redirect(w, r, "/admin/disciplines", http.StatusSeeOther)
		return
	}

	alerts.FlashSuccess(w, r, "Дисциплина обновлена")
	http.Redirect(w, r, "/admin/disciplines", http.StatusSeeOther)
}

func (h *Handler) DeleteDiscipline(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusNotFound)
		return
	}

	err = r.ParseForm()
	if err != nil {
		alerts.FlashError(w, r, "Ошибка обработки формы")
		log.Printf("Failed to parse form: %v", err)
		http.Redirect(w, r, "/admin/disciplines", http.StatusSeeOther)
		return
	}

	delDisc, err := h.db.GetDisciplineByID(id)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка чтения дисциплины из БД")
		log.Printf("Failed to get disc from db: %v", err)
		http.Redirect(w, r, "/admin/disciplines", http.StatusSeeOther)
		return
	}

	if delDisc.Name != r.FormValue("name") {
		alerts.FlashWarning(w, r, "Дисциплина не удалена: указано неверное имя дисциплины")
		http.Redirect(w, r, "/admin/disciplines", http.StatusSeeOther)
		return
	}

	err = h.db.DeleteDiscipline(delDisc)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка удаления дисциплины из БД")
		log.Printf("Failed to delete disc: %v", err)
		http.Redirect(w, r, "/admin/disciplines", http.StatusSeeOther)
		return
	}

	alerts.FlashSuccess(w, r, "Дисциплина удалена")
	http.Redirect(w, r, "/admin/disciplines", http.StatusSeeOther)
}
