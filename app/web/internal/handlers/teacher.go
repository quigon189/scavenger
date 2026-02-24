package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"scavenger/core/models"
	"scavenger/core/services"
	"scavenger/web/internal/services/session"
	"scavenger/web/internal/views/pages/teacher"
)

type TeacherHandler struct {
	services *services.Services
	session  *session.SessionService
}

func NewTeacherHandler(services *services.Services, session *session.SessionService) *TeacherHandler {
	return &TeacherHandler{
		services: services,
		session:  session,
	}
}

func (h *TeacherHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)

	disciplines, err := h.services.Data.GetDisciplinesByTeacher(r.Context(), session.User.ID)
	if err != nil {
		h.session.FlashError(w, r, "Ошибка при загрузке дисциплин", err)
		disciplines = []models.Discipline{}
	}

	BaseWithNavbar(w, r, "Мои дисциплины", teacher.Dashboard(session.User, disciplines))
}

func (h *TeacherHandler) CreateDisciplineForm(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)

	groups, err := h.services.Data.GetAllGroups(r.Context())
	if err != nil {
		h.session.FlashError(w, r, "Ошибка при загрузке групп", err)
		groups = []models.Group{}
	}

	periods, err := h.services.Data.GetAllPeriods(r.Context())
	if err != nil {
		h.session.FlashError(w, r, "Ошибка при загрузке периодов", err)
		periods = []models.Period{}
	}

	BaseWithNavbar(w, r, "Создание дисциплины", teacher.DisciplineForm(session.User, nil, groups, periods))
}

func (h *TeacherHandler) CreateDiscipline(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)

	name := r.FormValue("name")
	description := r.FormValue("description")
	groupID, _ := strconv.Atoi(r.FormValue("group_id"))
	periodID, _ := strconv.Atoi(r.FormValue("period_id"))

	discipline := &models.Discipline{
		Name:        name,
		Description: description,
		TeacherID:   session.User.ID,
		GroupID:     groupID,
		PeriodID:    periodID,
	}

	err := h.services.Data.CreateDiscipline(r.Context(), discipline)
	if err != nil {
		h.session.FlashError(w, r, "Ошибка при создании дисциплины", err)
		http.Redirect(w, r, "/teacher/disciplines/new", http.StatusSeeOther)
		return
	}

	h.session.FlashSuccess(w, r, "Дисциплина успешно создана")
	http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
}

func (h *TeacherHandler) EditDisciplineForm(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.session.FlashError(w, r, "Неверный ID дисциплины", err)
		http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
		return
	}

	discipline, err := h.services.Data.GetDiscipline(r.Context(), id)
	if err != nil || discipline == nil {
		h.session.FlashError(w, r, "Дисциплина не найдена", err)
		http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
		return
	}

	if discipline.TeacherID != session.User.ID {
		h.session.FlashError(w, r, "Нет прав на редактирование этой дисциплины", fmt.Errorf("access denied"))
		http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
		return
	}

	groups, err := h.services.Data.GetAllGroups(r.Context())
	if err != nil {
		h.session.FlashError(w, r, "Ошибка при загрузке групп", err)
		groups = []models.Group{}
	}

	periods, err := h.services.Data.GetAllPeriods(r.Context())
	if err != nil {
		h.session.FlashError(w, r, "Ошибка при загрузке периодов", err)
		periods = []models.Period{}
	}

	BaseWithNavbar(w, r, "Редактирование дисциплины", teacher.DisciplineForm(session.User, discipline, groups, periods))
}

func (h *TeacherHandler) EditDiscipline(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.session.FlashError(w, r, "Неверный ID дисциплины", err)
		http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	groupID, _ := strconv.Atoi(r.FormValue("group_id"))
	periodID, _ := strconv.Atoi(r.FormValue("period_id"))

	discipline := &models.Discipline{
		ID:          id,
		Name:        name,
		Description: description,
		TeacherID:   session.User.ID,
		GroupID:     groupID,
		PeriodID:    periodID,
	}

	err = h.services.Data.UpdateDiscipline(r.Context(), discipline)
	if err != nil {
		h.session.FlashError(w, r, "Ошибка при обновлении дисциплины", err)
		http.Redirect(w, r, "/teacher/disciplines/"+idStr+"/edit", http.StatusSeeOther)
		return
	}

	h.session.FlashSuccess(w, r, "Дисциплина обновлена")
	http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
}

func (h *TeacherHandler) DeleteDiscipline(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.session.FlashError(w, r, "Неверный ID дисциплины", err)
		http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
		return
	}

	err = h.services.Data.DeleteDiscipline(r.Context(), id)
	if err != nil {
		h.session.FlashError(w, r, "Ошибка при удалении дисциплины", err)
	} else {
		h.session.FlashSuccess(w, r, "Дисциплина удалена")
	}
	http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
}

func (h *TeacherHandler) LabsList(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)
	disciplineIDStr := r.PathValue("disciplineId")
	disciplineID, err := strconv.Atoi(disciplineIDStr)
	if err != nil {
		h.session.FlashError(w, r, "Неверный ID дисциплины", err)
		http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
		return
	}

	discipline, err := h.services.Data.GetDiscipline(r.Context(), disciplineID)
	if err != nil || discipline == nil {
		h.session.FlashError(w, r, "Дисциплина не найдена", err)
		http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
		return
	}
	if discipline.TeacherID != session.User.ID {
		h.session.FlashError(w, r, "Нет доступа к этой дисциплине", fmt.Errorf("access denied"))
		http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
		return
	}

	labs, err := h.services.Data.GetLabsByDiscipline(r.Context(), disciplineID)
	if err != nil {
		h.session.FlashError(w, r, "Ошибка при загрузке лабораторных работ", err)
		labs = []models.Lab{}
	}

	BaseWithNavbar(w, r, "Лабораторные работы", teacher.LabsList(session.User, discipline, labs))
}

func (h *TeacherHandler) CreateLabForm(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)
	disciplineIDStr := r.PathValue("disciplineId")
	disciplineID, err := strconv.Atoi(disciplineIDStr)
	if err != nil {
		h.session.FlashError(w, r, "Неверный ID дисциплины", err)
		http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
		return
	}

	discipline, err := h.services.Data.GetDiscipline(r.Context(), disciplineID)
	if err != nil || discipline == nil {
		h.session.FlashError(w, r, "Дисциплина не найдена", err)
		http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
		return
	}
	if discipline.TeacherID != session.User.ID {
		h.session.FlashError(w, r, "Нет доступа к этой дисциплине", fmt.Errorf("access denied"))
		http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
		return
	}

	BaseWithNavbar(w, r, "Создание лабораторной работы", teacher.LabForm(session.User, discipline, nil))
}

func (h *TeacherHandler) CreateLab(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)
	disciplineIDStr := r.PathValue("disciplineId")
	disciplineID, err := strconv.Atoi(disciplineIDStr)
	if err != nil {
		h.session.FlashError(w, r, "Неверный ID дисциплины", err)
		http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
		return
	}

	discipline, err := h.services.Data.GetDiscipline(r.Context(), disciplineID)
	if err != nil || discipline == nil {
		h.session.FlashError(w, r, "Дисциплина не найдена", err)
		http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
		return
	}
	if discipline.TeacherID != session.User.ID {
		h.session.FlashError(w, r, "Нет доступа к этой дисциплине", fmt.Errorf("access denied"))
		http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	mdContent := r.FormValue("md_content")
	deadlineStr := r.FormValue("deadline")

	var deadline time.Time
	if deadlineStr != "" {
		deadline, err = time.Parse("2006-01-02T15:04", deadlineStr)
		if err != nil {
			h.session.FlashError(w, r, "Неверный формат даты", err)
			http.Redirect(w, r, "/teacher/disciplines/"+disciplineIDStr+"/labs/new", http.StatusSeeOther)
			return
		}
	}

	lab := &models.Lab{
		DisciplineID: disciplineID,
		Name:         name,
		Description:  description,
		MDContent:    mdContent,
		Deadline:     deadline,
	}

	err = h.services.Data.CreateLab(r.Context(), lab)
	if err != nil {
		h.session.FlashError(w, r, "Ошибка при создании лабораторной работы", err)
		http.Redirect(w, r, "/teacher/disciplines/"+disciplineIDStr+"/labs/new", http.StatusSeeOther)
		return
	}

	h.session.FlashSuccess(w, r, "Лабораторная работа создана")
	http.Redirect(w, r, "/teacher/disciplines/"+disciplineIDStr+"/labs", http.StatusSeeOther)
}

func (h *TeacherHandler) EditLabForm(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)
	labIDStr := r.PathValue("labId")
	labID, err := strconv.Atoi(labIDStr)
	if err != nil {
		h.session.FlashError(w, r, "Неверный ID лабораторной", err)
		http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
		return
	}

	lab, err := h.services.Data.GetLab(r.Context(), labID)
	if err != nil || lab == nil {
		h.session.FlashError(w, r, "Лабораторная работа не найдена", err)
		http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
		return
	}
	if lab.Discipline.TeacherID != session.User.ID {
		h.session.FlashError(w, r, "Нет доступа к этой лабораторной", fmt.Errorf("access denied"))
		http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
		return
	}

	discipline := &lab.Discipline

	BaseWithNavbar(w, r, "Редактирование лабораторной работы", teacher.LabForm(session.User, discipline, lab))
}

func (h *TeacherHandler) EditLab(w http.ResponseWriter, r *http.Request) {
	session := getSession(r)
	labIDStr := r.PathValue("labId")
	labID, err := strconv.Atoi(labIDStr)
	if err != nil {
		h.session.FlashError(w, r, "Неверный ID лабораторной", err)
		http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
		return
	}

	lab, err := h.services.Data.GetLab(r.Context(), labID)
	if err != nil || lab == nil {
		h.session.FlashError(w, r, "Лабораторная работа не найдена", err)
		http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
		return
	}
	if lab.Discipline.TeacherID != session.User.ID {
		h.session.FlashError(w, r, "Нет доступа к этой лабораторной", fmt.Errorf("access denied"))
		http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	mdContent := r.FormValue("md_content")
	deadlineStr := r.FormValue("deadline")

	var deadline time.Time
	if deadlineStr != "" {
		deadline, err = time.Parse("2006-01-02T15:04", deadlineStr)
		if err != nil {
			h.session.FlashError(w, r, "Неверный формат даты", err)
			http.Redirect(w, r, "/teacher/labs/"+labIDStr+"/edit", http.StatusSeeOther)
			return
		}
	}

	lab.Name = name
	lab.Description = description
	lab.MDContent = mdContent
	lab.Deadline = deadline

	err = h.services.Data.UpdateLab(r.Context(), lab)
	if err != nil {
		h.session.FlashError(w, r, "Ошибка при обновлении лабораторной работы", err)
		http.Redirect(w, r, "/teacher/labs/"+labIDStr+"/edit", http.StatusSeeOther)
		return
	}

	h.session.FlashSuccess(w, r, "Лабораторная работа обновлена")
	http.Redirect(w, r, "/teacher/disciplines/"+strconv.Itoa(lab.DisciplineID)+"/labs", http.StatusSeeOther)
}

func (h *TeacherHandler) DeleteLab(w http.ResponseWriter, r *http.Request) {
	labIDStr := r.PathValue("labId")
	labID, err := strconv.Atoi(labIDStr)
	if err != nil {
		h.session.FlashError(w, r, "Неверный ID лабораторной", err)
		http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
		return
	}

	// Получим дисциплину для редиректа
	lab, err := h.services.Data.GetLab(r.Context(), labID)
	if err != nil || lab == nil {
		h.session.FlashError(w, r, "Лабораторная работа не найдена", err)
		http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
		return
	}
	disciplineID := lab.DisciplineID

	err = h.services.Data.DeleteLab(r.Context(), labID)
	if err != nil {
		h.session.FlashError(w, r, "Ошибка при удалении лабораторной работы", err)
	} else {
		h.session.FlashSuccess(w, r, "Лабораторная работа удалена")
	}
	http.Redirect(w, r, "/teacher/disciplines/"+strconv.Itoa(disciplineID)+"/labs", http.StatusSeeOther)
}
