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

	//reports, totalCount, err := getFilteredReports(filterParams)
	reports, err := h.db.GetAllReports()
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

func (h *Handler) GradeModalHandler(w http.ResponseWriter, r *http.Request) {
	reportID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Неверный ID отчета", http.StatusBadRequest)
		return
	}

	report, err := h.db.GetLabReportByID(reportID)
	if err != nil {
		http.Error(w, "Отчет не найден", http.StatusNotFound)
		return
	}

	views.GradeModal(*report).Render(r.Context(), w)
}
