package handlers

import (
	"log"
	"net/http"
	"scavenger/internal/alerts"
	"scavenger/internal/auth"
	"scavenger/internal/models"
	"scavenger/views"
	"sort"
	"strconv"
)

func (h *Handler) StudentDashboard(w http.ResponseWriter, r *http.Request) {
	stud := models.GetUserFromContext(r.Context())
	disciplines, err := h.db.GetDisciplinesByGroupID(stud.GroupID)
	if err != nil {
		alerts.FlashWarning(w,r,"Дисциплины не загружены")
		disciplines = []models.Discipline{}
	}

	for i, discipline := range disciplines {
		labs,err  := h.db.GetDisciplineLabs(discipline.ID)
		if err != nil {
			continue
		}
		disciplines[i].Labs = append(disciplines[i].Labs, labs...)
	}

	reports := []models.LabReport{}
	for _, disc := range disciplines {
		for _, lab := range disc.Labs {
			labID, err := strconv.Atoi(lab.ID)
			if err != nil {
				continue
			}
			report, err := h.db.GetLabReport(stud.ID, labID)
			if err != nil {
				continue
			}

			report.Discipline = disc
			report.Lab = lab

			reports = append(reports, *report)
		}
	}

	sort.Slice(reports, func(i, j int) bool {
		return reports[i].UploadedAt.After(reports[j].UploadedAt)
	})

	if len(reports) > 5 {
		reports = reports[:5]
	}

	views.StudentDashboard(disciplines, reports).Render(r.Context(), w)
}

func (h *Handler) DisciplinePage(w http.ResponseWriter, r *http.Request) {
	id,err := strconv.Atoi(r.PathValue("id"))

	disc, err := h.db.GetDisciplineWithLabs(id)
	if err != nil {
		alerts.FlashError(w,r,"Ошибка чтения дисциплины из БД")
		log.Printf("Failed to get discipline: %v", err)
		http.Redirect(w,r,"/",http.StatusSeeOther)
		return
	}

	student := models.GetUserFromContext(r.Context())

	for i, lab := range disc.Labs {
		labID, err := strconv.Atoi(lab.ID)
		if err != nil {
			continue
		}
		report, err := h.db.GetLabReport(student.ID, labID)
		if err != nil {
			continue
		}

		disc.Labs[i].Reports = append(disc.Labs[i].Reports, *report)
	}

	views.DisciplinePage(disc).Render(r.Context(), w)
}

func (h *Handler) StudentRegisterPage(w http.ResponseWriter, r *http.Request) {
	groups, err := h.db.GetAllGroups()
    if err != nil {
		alerts.FlashError(w,r,"Ошибка получения групп")
        log.Printf("Failed to get groups: %v", err)
        groups = []models.Group{}
    }
    
	views.StudentRegistrationPage(groups).Render(r.Context(), w)
}

func (h *Handler) RegisterStudent(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        alerts.FlashError(w, r, "Ошибка обработки формы")
        http.Redirect(w, r, "/register", http.StatusSeeOther)
        return
    }

    username := r.FormValue("username")
    name := r.FormValue("name")
    password := r.FormValue("password")
    groupIDStr := r.FormValue("group_id")

    if username == "" || name == "" || password == "" || groupIDStr == "" {
        alerts.FlashError(w, r, "Все поля обязательны для заполнения")
        http.Redirect(w, r, "/register", http.StatusSeeOther)
        return
    }

    if len(password) < 6 {
        alerts.FlashError(w, r, "Пароль должен содержать минимум 6 символов")
        http.Redirect(w, r, "/register", http.StatusSeeOther)
        return
    }

    groupID, err := strconv.Atoi(groupIDStr)
    if err != nil {
        alerts.FlashError(w, r, "Неверный формат группы")
        http.Redirect(w, r, "/register", http.StatusSeeOther)
        return
    }

    group, err := h.db.GetGroupByID(groupID)
    if err != nil {
        alerts.FlashError(w, r, "Выбранная группа не найдена")
        http.Redirect(w, r, "/register", http.StatusSeeOther)
        return
    }

    existingUser, err := h.db.GetUserByUsername(username)
    if err == nil && existingUser != nil {
        alerts.FlashError(w, r, "Пользователь с таким логином уже существует")
        http.Redirect(w, r, "/register", http.StatusSeeOther)
        return
    }

    newStudent, err := auth.RegisterUser(
        username,
        name,
        password,
        string(models.StudentRole),
    )
    if err != nil {
        log.Printf("Ошибка регистрации студента: %v", err)
        alerts.FlashError(w, r, "Ошибка при создании учетной записи")
        http.Redirect(w, r, "/register", http.StatusSeeOther)
        return
    }

    newStudent.GroupID = group.ID
    newStudent.GroupName = group.Name

    if err := h.db.CreateStudent(newStudent); err != nil {
        log.Printf("Ошибка сохранения студента в БД: %v", err)
        alerts.FlashError(w, r, "Ошибка при сохранении данных студента")
        http.Redirect(w, r, "/register", http.StatusSeeOther)
        return
    }

    log.Printf("Зарегистрирован новый студент: %s (%s), группа: %s", 
        newStudent.Name, newStudent.Username, newStudent.GroupName)

    alerts.FlashSuccess(w, r, "Регистрация успешно завершена! Теперь вы можете войти в систему.")
    http.Redirect(w, r, "/", http.StatusSeeOther)
}
