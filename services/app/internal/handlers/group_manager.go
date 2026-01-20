package handlers

import (
	"log"
	"net/http"
	"strconv"

	"scavenger/internal/alerts"
	"scavenger/internal/models"
	"scavenger/views"
)

func (h *Handler) GroupManager(w http.ResponseWriter, r *http.Request) {
	groups, err := h.db.GetAllGroupsWithStudentsAndDisciplines()
	if err != nil {
		alerts.FlashWarning(w, r, "Группы не получены")
		log.Printf("Failed to get groups: %v", err)
		groups = []models.Group{}
	}

	disciplines, err := h.db.GetDisciplinesWithoutGroup()
	if err != nil {
		alerts.FlashError(w, r, "Дисциплины не получены")
		log.Printf("Failed to get disciplines: %v", err)
		disciplines = []models.Discipline{}
	}

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

	discs := r.Form["existing_disciplines"]
	if len(discs) == 0 {
		alerts.FlashSuccess(w, r, "Группа "+newGroup.Name+" создана, дисциплины не выбраны")
		http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
		return
	}

	for _, disc := range discs {
		discID, err := strconv.Atoi(disc)
		if err != nil {
			alerts.FlashError(w, r, "Ошибка ID дисциплины")
			log.Printf("Failed to conv discipline id: %v", err)
			continue
		}
		editDisc, err := h.db.GetDisciplineByID(discID)
		if err != nil {
			alerts.FlashError(w, r, "Ошибка чтения дисциплины из БД")
			log.Printf("Failed to get discipline: %v", err)
			continue
		}

		editDisc.GroupID = &newGroup.ID

		err = h.db.UpdateDiscipline(editDisc)
		if err != nil {
			alerts.FlashError(w, r, "Ошибка обновления дисциплины")
			log.Printf("Failed to update discipline: %v", err)
			continue
		}
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

func (h *Handler) EditGroup(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusNotFound)
		return
	}

	err = r.ParseForm()
	if err != nil {
		alerts.FlashError(w, r, "Ошибка чтения формы")
		log.Printf("Failed to parse form: %v", err)
		http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
		return
	}

	editGroup, err := h.db.GetGroupByID(id)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка чтения группы из БД")
		log.Printf("Failed to get group: %v", err)
		http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
		return
	}

	editGroup.Name = r.FormValue("name")

	err = h.db.UpdateGroup(editGroup)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка обнавления группы")
		log.Printf("Failed to update group: %v", err)
		http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
		return
	}

	alerts.FlashSuccess(w, r, "Группа обновлена")
	http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
}

func (h *Handler) AddDiscToGroup(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusNotFound)
		return
	}

	err = r.ParseForm()
	if err != nil {
		alerts.FlashError(w, r, "Ошибка чтения формы")
		log.Printf("Failed to parse form: %v", err)
		http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
		return
	}

	editGroup, err := h.db.GetGroupByID(id)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка чтения группы из БД")
		log.Printf("Failed to get group: %v", err)
		http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
		return
	}

	discs := r.Form["existing_disciplines"]
	if len(discs) == 0 {
		alerts.FlashWarning(w, r, "Дисциплины не выбраны")
		http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
		return
	}

	for _, disc := range discs {
		discID, err := strconv.Atoi(disc)
		if err != nil {
			alerts.FlashError(w, r, "Ошибка ID дисциплины")
			log.Printf("Failed to conv discipline id: %v", err)
			continue
		}
		editDisc, err := h.db.GetDisciplineByID(discID)
		if err != nil {
			alerts.FlashError(w, r, "Ошибка чтения дисциплины из БД")
			log.Printf("Failed to get discipline: %v", err)
			continue
		}

		editDisc.GroupID = &editGroup.ID

		err = h.db.UpdateDiscipline(editDisc)
		if err != nil {
			alerts.FlashError(w, r, "Ошибка обновления дисциплины")
			log.Printf("Failed to update discipline: %v", err)
			continue
		}
	}

	alerts.FlashSuccess(w, r, "Дисциплины привязаны")
	http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
}

func (h *Handler) RemoveDiscFromGroup(w http.ResponseWriter, r *http.Request) {
	groupID, err := strconv.Atoi(r.PathValue("groupID"))
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusNotFound)
		return
	}
	discID, err := strconv.Atoi(r.PathValue("discID"))
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusNotFound)
		return
	}

	err = r.ParseForm()
	if err != nil {
		alerts.FlashError(w, r, "Ошибка чтения формы")
		log.Printf("Failed to parse form: %v", err)
		http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
		return
	}

	disc, err := h.db.GetDisciplineByID(discID)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка чтения дисциплины из БД")
		log.Printf("Failed to get discipline: %v", err)
		http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
		return
	}

	if groupID == *disc.GroupID {
		disc.GroupID = nil
		err = h.db.UpdateDiscipline(disc)
		if err != nil {
			alerts.FlashError(w, r, "Ошибка обновления дисциплины")
			log.Printf("Failed to update discipline: %v", err)
			http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
			return
		}
	}

	alerts.FlashSuccess(w, r, "Дисциплина отвязана")
	http.Redirect(w, r, "/admin/groups", http.StatusSeeOther)
}
