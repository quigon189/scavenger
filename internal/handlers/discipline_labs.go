package handlers

import (
	"log"
	"net/http"
	"path/filepath"
	"scavenger/internal/alerts"
	filestorage "scavenger/internal/file_storage"
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

	lab.DisciplineID = id

	storedFile, err := h.fs.SaveLabFile(
		filestorage.Markdowm,
		lab,
		MDFile,
		header,
	)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка сохранения файла")
		log.Printf("Failed to save file: %v", err)
		http.Redirect(w, r, "/admin/disciplines/"+r.PathValue("id"), http.StatusSeeOther)
		return
	}

	lab.MDFile.Path = storedFile.Path
	lab.MDFile.URL = storedFile.URL
	lab.MDFile.Filename = storedFile.Filename
	lab.MDFile.Size = storedFile.Size

	err = h.db.AddStoredFile(&lab.MDFile)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка сохранения файла в БД")
		log.Printf("Failed to save file to db: %v", err)
		http.Redirect(w, r, "/admin/disciplines/"+r.PathValue("id"), http.StatusSeeOther)
		return
	}

	lab.MDFileID = lab.MDFile.ID

	otherFiles := r.MultipartForm.File["files"]
	storedFiles := []models.StoredFile{}
	for _, fileHeader := range otherFiles {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		storedFile, err := h.fs.SaveLabFile(
			filestorage.LabMaterial,
			lab,
			file,
			fileHeader,
		)
		if err != nil {
			alerts.FlashError(w, r, "Ошибка сохранения файла")
			log.Printf("Failed to save file: %v", err)
			http.Redirect(w, r, "/admin/disciplines/"+r.PathValue("id"), http.StatusSeeOther)
			return

		}

		sFile := models.StoredFile{
			Path:     storedFile.Path,
			URL:      storedFile.URL,
			Filename: storedFile.Filename,
			Size:     storedFile.Size,
		}

		err = h.db.AddStoredFile(&sFile)
		if err == nil {
			storedFiles = append(storedFiles, sFile)
		}
	}

	lab.StoredFiles = append(lab.StoredFiles, storedFiles...)

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
			storedFile, err := h.fs.SaveLabFile(
				filestorage.Markdowm,
				lab,
				MDFile,
				header,
			)
			if err != nil {
				alerts.FlashError(w, r, "Ошибка сохранения файла")
				log.Printf("Failed to save file: %v", err)
				http.Redirect(w, r, "/admin/disciplines/"+r.PathValue("id"), http.StatusSeeOther)
				return
			}

			lab.MDFile.Path = storedFile.Path
			lab.MDFile.URL = storedFile.URL
			lab.MDFile.Filename = storedFile.Filename
			lab.MDFile.Size = storedFile.Size

			err = h.db.AddStoredFile(&lab.MDFile)
			if err != nil {
				alerts.FlashError(w, r, "Ошибка сохранения файла в БД")
				log.Printf("Failed to save file to db: %v", err)
				http.Redirect(w, r, "/admin/disciplines/"+r.PathValue("id"), http.StatusSeeOther)
				return
			}

			lab.MDFileID = lab.MDFile.ID
		}
	}

	var newFiles []models.StoredFile

	files := r.MultipartForm.File["files"]
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		storedFile, err := h.fs.SaveLabFile(
			filestorage.LabMaterial,
			lab,
			file,
			fileHeader,
		)
		if err != nil {
			alerts.FlashError(w, r, "Ошибка сохранения файла")
			log.Printf("Failed to save file: %v", err)
			http.Redirect(w, r, "/admin/disciplines/"+r.PathValue("id"), http.StatusSeeOther)
			return
		}

		sFile := models.StoredFile{
			Path:     storedFile.Path,
			URL:      storedFile.URL,
			Filename: storedFile.Filename,
			Size:     storedFile.Size,
		}

		err = h.db.AddStoredFile(&sFile)
		if err == nil {
			newFiles = append(newFiles, sFile)
		}
	}
	if len(newFiles) > 0 {
		h.db.AddLabFiles(labID, newFiles)
	}

	delFilesID := r.Form["remove_file"]
	if len(delFilesID) > 0 {
		delStoredFiles := []models.StoredFile{}
		for _, idS := range delFilesID {
			id, err := strconv.Atoi(idS)
			if err != nil {
				continue
			}
			delStoredFiles = append(delStoredFiles, models.StoredFile{
				ID: id,
			})
		}
		h.db.RemoveLabFiles(labID, delStoredFiles)
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
