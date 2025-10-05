package handlers

import (
	"log"
	"net/http"
	"strconv"

	"scavenger/internal/alerts"
	"scavenger/internal/models"
	"scavenger/views"
)

func (h *Handler) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	groups, err := h.db.GetAllGroupsWithDisciplines()
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

	disciplines, err := h.db.GetDisciplines()
	if err != nil {
		alerts.FlashWarning(w,r,"Дисциплины не загруженны")
		log.Printf("Failed to get disciplines: %v", err)
		disciplines = []models.Discipline{}
	}

	views.AdminDashboard(&stats, reports, disciplines).Render(r.Context(), w)
}

func (h *Handler) GroupManager(w http.ResponseWriter, r *http.Request) {
	groups, err := h.db.GetAllGroupsWithStudents()
	if err != nil {
		alerts.FlashError(w, r, "Ошибка получения групп")
		log.Printf("Failed to get groups: %v", err)
		groups = []models.Group{}
	}

	disciplines := []models.Discipline{}

	views.GroupsManager(groups, disciplines).Render(r.Context(), w)
}

func (h *Handler) AddGroup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		alerts.FlashError(w, r, "Ошибка обработки формы")
		log.Printf("Failed to parse form: %v", err)
		http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
	}

	groupName := r.FormValue("name")

	groups, err := h.db.GetAllGroups()
	if err != nil {
		alerts.FlashError(w, r, "Ошибка при получении групп")
		log.Printf("Failed to get groups: %v", err)
		http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
	}

	for _, g := range groups {
		if g.Name == groupName {
			alerts.FlashError(w, r, "Группа с таким именем уже существует")
			http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
		}
	}

	newGroup := &models.Group{
		Name: groupName,
	}

	if err := h.db.CreateGroup(newGroup); err != nil {
		alerts.FlashError(w, r, "Ошибка при создании группы")
		log.Printf("Failed to create group: %v", err)
		http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
	}

	alerts.FlashSuccess(w, r, "Группа "+newGroup.Name+" создана")

	http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
}

func (h *Handler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		alerts.FlashError(w, r, "Ошибка чтения формы")
		log.Printf("Failed to parse form: %v", err)
		http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
	}
	groupName := r.FormValue("name")

	groupID, err := strconv.Atoi(r.PathValue("groupID"))
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusNotFound)
	}

	group, err := h.db.GetGroupByID(groupID)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка при получении группы")
		log.Printf("Failed to get group: %v", err)
		http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
	}

	if groupName != group.Name {
		alerts.FlashError(w, r, "Введено неверное имя группы")
		http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
	}

	if err := h.db.DeleteGroupByID(group.ID); err != nil {
		alerts.FlashError(w, r, "Ошибка при удалении группы")
		log.Printf("Failed to delete group: %v", err)
		http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
	}

	alerts.FlashSuccess(w, r, "Группа удалена")
	http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
}
