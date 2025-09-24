package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"scavenger/internal/config"
	"scavenger/internal/models"
	"scavenger/views"
	"time"
)

func (h *Handler) StudentDashboard(w http.ResponseWriter, r *http.Request) {
	username := h.authService.GetUsername(r)
	group := h.authService.GetGroup(r)

	disciplines := h.cfg.GetGroupDisciplines(group)
	reports := h.cfg.GetStudentReports(username)

	log.Printf("Render stud dashboard: group %s, disciplines %+v, reports %+v", group, disciplines, reports)
	views.StudentDashboard(disciplines, group, reports).Render(r.Context(), w)
}

func (h *Handler) DisciplinePage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	group := h.authService.GetGroup(r)
	disc := h.cfg.GetDiscepline(id)

	if disc.Name == "" {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
	}

	views.DisciplinePage(*disc, group).Render(r.Context(), w)
}

func (h *Handler) DownloadLabs(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	log.Printf("download path: %s", path)

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
	discipline := r.FormValue("discipline")
	comment := r.FormValue("comment")

	file, header, err := r.FormFile("report_file")
	if err != nil {
		http.Error(w, "Ошибка получения файла: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	timestamp := time.Now().Format("20060101-154015")
	newFilename := fmt.Sprintf("%s-lab_id_%s-%s.%s",
		timestamp,
		labID,
		username,
		ext,
	)

	uploadPath := filepath.Join(h.cfg.Server.UploadPath, "reports", newFilename)

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

	report := models.LabReport{
		ID:         timestamp,
		Student:    username,
		Group:      group,
		Discipline: discipline,
		LabName:    h.cfg.GetLab(labID).Name,
		Path:       uploadPath,
		Comment:    comment,
		UploadetAt: time.Now(),
		Status:     "submitted",
	}

	h.cfg.LabReports = append(h.cfg.LabReports, report)

	config.SaveConfig(h.cfg)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
