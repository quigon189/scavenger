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
		http.Redirect(w, r, "/404", http.StatusNotFound)
		return
	}

	discipline, err := h.db.GetDisciplineByID(id)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка чтения дисциплины из БД")
		log.Printf("Failed to get disciplint from db: %v", err)
		http.Redirect(w, r, "/admin/disciplines", http.StatusInternalServerError)
		return
	}

	views.DisciplineLabs(discipline).Render(r.Context(), w)
}

func (h *Handler) AddDisciplineLabs(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusNotFound)
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

	log.Printf("Save file: %s", header.Filename)
	lab.MDPath = header.Filename

	pdfFiles := r.MultipartForm.File["pdf_files"]
	for _, fileHeader := range pdfFiles {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		lab.PDFPath = append(lab.PDFPath, fileHeader.Filename)
	}

	lab.DisciplineID = id

	err = h.db.AddDisciplineLab(lab)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка сохраниния записи в БД")
		log.Printf("Failed to add lab: %v", err)
		http.Redirect(w, r, "/admin/disciplines/"+r.PathValue("id"), http.StatusSeeOther)
		return
	}

	log.Printf("Save Lab to disc %v: %+v", id, lab)
	alerts.FlashSuccess(w, r, "Работа добавлена")
	http.Redirect(w, r, "/admin/disciplines/"+r.PathValue("id"), http.StatusSeeOther)
}
