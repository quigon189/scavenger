package handlers

import (
	"fmt"
	"log"
	"net/http"
	"scavenger/internal/alerts"
	"scavenger/internal/models"
	"scavenger/views"
	"strconv"
)

func (h *Handler) ReportsPage(w http.ResponseWriter, r *http.Request) {
	disciplines, err := h.db.GetDisciplines()
	if err != nil {
		alerts.FlashWarning(w, r, "Дисциплины не загружены")
		log.Printf("Ошибка получения дисциплин: %v", err)
		disciplines = []models.Discipline{}
	}

	for i, disc := range disciplines {
		labs, err := h.db.GetDisciplineLabs(disc.ID)
		if err != nil {
			continue
		}

		disciplines[i].Labs = append(disciplines[i].Labs, labs...)
	}

	var filterParams models.ReportFilterParams

	filterParams.Parse(r)

	reports, err := h.db.GetFilteredReports(filterParams)
	if err != nil {
		log.Printf("Ошибка получения отчетов: %v", err)
		reports = []models.LabReport{}
	}

	totalCount := len(reports)

	filterParams.TotalPages = (totalCount + filterParams.PageSize - 1) / filterParams.PageSize

	views.ReportsPage(disciplines, reports, filterParams).Render(r.Context(), w)
}

func (h *Handler) ReportsTable(w http.ResponseWriter, r *http.Request) {
	var filterParams models.ReportFilterParams

	filterParams.Parse(r)

	reports, err := h.db.GetFilteredReports(filterParams)
	if err != nil {
		log.Printf("Ошибка получения отчетов: %v", err)
		reports = []models.LabReport{}
	}

	totalCount := len(reports)

	filterParams.TotalPages = (totalCount + filterParams.PageSize - 1) / filterParams.PageSize

	views.ReportsTable(reports, filterParams).Render(r.Context(), w)
}

func (h *Handler) ReportsLabsByDiscipline(w http.ResponseWriter, r *http.Request) {
	disciplineIDStr := r.URL.Query().Get("discipline_id")
	if disciplineIDStr == "" {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<option value="">Все работы</option>`)
		return
	}

	disciplineID, err := strconv.Atoi(disciplineIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	labs, err := h.db.GetDisciplineLabs(disciplineID)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<option value="">Нет доступных работ</option>`)
		return
	}

	fmt.Fprint(w, `<option value="">Все работы</option>`)
	for _, lab := range labs {
		fmt.Fprintf(w, `<option value="%s">%s</option>`, lab.ID, lab.Name)
	}
}

func (h *Handler) ReportReviewPage(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	report, err := h.db.GetLabReportByID(id)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка чтения отчета")
		log.Printf("Failed to get LabReport: %v", err)
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	student, err := h.db.GetStudentByID(report.StudentID)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка чтения студента")
		log.Printf("Failed to get student: %v", err)
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}
	report.Student = *student

	discipline, err := h.db.GetDisciplineByID(report.DisciplineID)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка чтения дисциплины")
		log.Printf("Failed to get discipline: %v", err)
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}
	report.Discipline = *discipline

	lab, err := h.db.GetLabByID(report.LabID)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка чтения работы")
		log.Printf("Failed to get lab: %v", err)
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}
	report.Lab = *lab

	views.ReportReviewPage(*report).Render(r.Context(), w)
}

func (h *Handler) GradeReport(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		alerts.FlashError(w, r, "Ошибка при обработке формы")
		log.Printf("Failed to parse form: %v", err)
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	grade, _ := strconv.Atoi(r.FormValue("grade"))
	teacherNote := r.FormValue("teacher_note")

	if grade < 2 || grade > 5 {
		alerts.FlashError(w, r, "Неверная оценка")
		log.Printf("Failed to validate grade")
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}
	
	report, err := h.db.GetLabReportByID(id)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка чтения отчета")
		log.Printf("Failed to get report: %v", err)
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	report.Grade = grade
	report.Status = "graded"
	report.TeacherNote = teacherNote

	err = h.db.UpdateReportGrade(report)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка при сохранении оценки")
		log.Printf("Failed to update report: %v", err)
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		return
	}

	alerts.FlashSuccess(w, r, "Оценка успешно сохранена")
	http.Redirect(w, r, "/admin/reports", http.StatusSeeOther)
}
