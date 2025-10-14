package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"scavenger/internal/alerts"
	"scavenger/internal/config"
	"scavenger/internal/models"
	"scavenger/views"
	"sort"
	"strconv"
	"time"
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

func (h *Handler) DownloadLabs(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")

	filePath := filepath.Join(h.cfg.Server.UploadPath, path)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
	w.Header().Set("Content-Type", "application/octet-stream")

	http.ServeFile(w, r, filePath)
}

func (h *Handler) UploadReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	username := h.authService.GetUsername(r)
	group := h.authService.GetGroup(r)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Ошибка загрузки файла: "+err.Error(), http.StatusBadRequest)
		return
	}

	labID := r.FormValue("lab_id")

	file, header, err := r.FormFile("report_file")
	if err != nil {
		http.Error(w, "Ошибка получения файла: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	timestamp := time.Now().Format("20060101-154015")
	newFilename := fmt.Sprintf("%s.%s",
		timestamp,
		ext,
	)

	uploadPath := filepath.Join(h.cfg.Server.UploadPath, "reports", group, username, labID, newFilename)

	os.MkdirAll(filepath.Dir(uploadPath), 0755)

	outFile, err := os.Create(uploadPath)
	if err != nil {
		http.Error(w, "Ошибка сохранения файла: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Ошибка копирования файла: "+err.Error(), http.StatusInternalServerError)
		return
	}

	report := models.LabReport{}

	log.Printf("report: %v", report)

	err = config.SaveConfig(h.cfg)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка при сохранении")
		log.Printf("Error to save config: %v", err)
		return
	}

	alerts.FlashSuccess(w, r, "Отчет загружен")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
