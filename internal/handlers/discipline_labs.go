package handlers

import (
	"log"
	"net/http"
	"path/filepath"
	"scavenger/internal/alerts"
	"scavenger/internal/models"
	"scavenger/views"
	"strconv"
	"time"
)

func (h *Handler) DisciplineLabs(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	discipline, err := h.db.GetDisciplineWithLabs(id)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка чтения дисциплины из БД")
		log.Printf("Failed to get discipline from db: %v", err)
		http.Redirect(w, r, "/admin/disciplines", http.StatusSeeOther)
		return
	}

	views.DisciplineLabs(discipline).Render(r.Context(), w)
}

func (h *Handler) AddDisciplineLabs(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка чтения формы")
		log.Printf("Failed to parse form: %v", err)
		http.Redirect(w, r, "/admin/disciplines/"+r.PathValue("id"), http.StatusSeeOther)
		return
	}

	lab := &models.Lab{}
	lab.Name = r.Form.Get("name")
	date, err := time.Parse("2006-01-02", r.Form.Get("deadline"))
	if err != nil {
		alerts.FlashError(w, r, "Ошибка чтения даты")
		log.Printf("Failed to parse date: %v", err)
		http.Redirect(w, r, "/admin/disciplines/"+r.PathValue("id"), http.StatusSeeOther)
		return
	}
	lab.Deadline = date
	lab.Description = r.Form.Get("description")

	MDFile, header, err := r.FormFile("md_file")
	if err != nil {
		alerts.FlashError(w, r, "Ошибка получения файла")
		log.Printf("Failed to download file: %v", err)
		http.Redirect(w, r, "/admin/disciplines/"+r.PathValue("id"), http.StatusSeeOther)
		return
	}
	defer MDFile.Close()

	ext := filepath.Ext(header.Filename)
	if ext != ".md" {
		alerts.FlashError(w, r, "Неверное расширение md файла")
		http.Redirect(w, r, "/admin/disciplines/"+r.PathValue("id"), http.StatusSeeOther)
		return
	}

	lab.MDPath = header.Filename

	pdfFiles := r.MultipartForm.File["files"]
	for _, fileHeader := range pdfFiles {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		lab.FilesPath = append(lab.FilesPath, fileHeader.Filename)
	}

	lab.DisciplineID = id

	err = h.db.AddDisciplineLab(lab)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка сохраниния записи в БД")
		log.Printf("Failed to add lab: %v", err)
		http.Redirect(w, r, "/admin/disciplines/"+r.PathValue("id"), http.StatusSeeOther)
		return
	}

	alerts.FlashSuccess(w, r, "Работа добавлена")
	http.Redirect(w, r, "/admin/disciplines/"+r.PathValue("id"), http.StatusSeeOther)
}

func (h *Handler) EditDisciplineLab(w http.ResponseWriter, r *http.Request) {
	labID, err := strconv.Atoi(r.PathValue("labID"))
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	lab, err := h.db.GetLabByID(labID)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка чтения работы")
		log.Printf("Failed to get lab: %v", err)
		http.Redirect(w, r, "/admin/disciplines/"+r.PathValue("discID"), http.StatusSeeOther)
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка чтения формы")
		log.Printf("Failed to parse form: %v", err)
		http.Redirect(w, r, "/admin/disciplines/"+r.PathValue("discID"), http.StatusSeeOther)
		return
	}

	lab.Name = r.Form.Get("name")
	date, err := time.Parse("2006-01-02", r.Form.Get("deadline"))
	if err != nil {
		alerts.FlashError(w, r, "Ошибка чтения даты")
		log.Printf("Failed to parse date: %v", err)
		http.Redirect(w, r, "/admin/disciplines/"+r.PathValue("discID"), http.StatusSeeOther)
		return
	}
	lab.Deadline = date
	lab.Description = r.Form.Get("description")

	MDFile, header, err := r.FormFile("md_file")
	if err != nil && err.Error() != "http: no such file" {
		alerts.FlashError(w, r, "Ошибка получения файла")
		log.Printf("Failed to download file: %v", err)
		http.Redirect(w, r, "/admin/disciplines/"+r.PathValue("discID"), http.StatusSeeOther)
		return
	}
	if err == nil || err.Error() != "http: no such file" {
		defer MDFile.Close()

		if header.Filename != "" {
			ext := filepath.Ext(header.Filename)
			if ext != ".md" {
				alerts.FlashError(w, r, "Неверное расширение md файла")
				http.Redirect(w, r, "/admin/disciplines/"+r.PathValue("discID"), http.StatusSeeOther)
				return
			}
			lab.MDPath = header.Filename
		}
	}

	var newFiles []string

	files := r.MultipartForm.File["files"]
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		lab.FilesPath = append(lab.FilesPath, fileHeader.Filename)
		newFiles = append(newFiles, fileHeader.Filename)
	}
	if len(newFiles) > 0 {
		h.db.AddLabFiles(labID, newFiles)
	}

	delFiles := r.Form["remove_file"]
	if len(delFiles) > 0 {
		h.db.RemoveLabFiles(labID, delFiles)
	}

	err = h.db.UpdateLab(lab)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка сохраниния записи в БД")
		log.Printf("Failed to add lab: %v", err)
		http.Redirect(w, r, "/admin/disciplines/"+r.PathValue("discID"), http.StatusSeeOther)
		return
	}

	alerts.FlashSuccess(w, r, "Работа обновлена")
	http.Redirect(w, r, "/admin/disciplines/"+r.PathValue("discID"), http.StatusSeeOther)
}

func (h *Handler) DeleteDisciplineLab(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("labID"))
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	err = r.ParseForm()
	if err != nil {
		alerts.FlashError(w, r, "Ошибка чтения формы")
		log.Printf("Failed to parse form: %v", err)
		http.Redirect(w, r, "/admin/disciplines/"+r.PathValue("discID"), http.StatusSeeOther)
		return
	}

	delLab, err := h.db.GetLabByID(id)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка чтения лабораторной работы")
		log.Printf("Failed to get lab: %v", err)
		http.Redirect(w, r, "/admin/disciplines/"+r.PathValue("discID"), http.StatusSeeOther)
		return
	}

	if delLab.Name != r.FormValue("name") {
		alerts.FlashWarning(w, r, "Указано неверное имя")
		http.Redirect(w, r, "/admin/disciplines/"+r.PathValue("discID"), http.StatusSeeOther)
		return
	}

	err = h.db.DeleteLab(id)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка удаления лабораторной работы")
		log.Printf("Failed to remove lab: %v", err)
		http.Redirect(w, r, "/admin/disciplines/"+r.PathValue("discID"), http.StatusSeeOther)
		return
	}

	alerts.FlashSuccess(w, r, "Работа удалена")
	http.Redirect(w, r, "/admin/disciplines/"+r.PathValue("discID"), http.StatusSeeOther)
}
