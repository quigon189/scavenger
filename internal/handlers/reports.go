package handlers

import (
	"log"
	"net/http"
	"scavenger/internal/alerts"
	"scavenger/internal/models"
	"scavenger/views"
	"strconv"
	"time"
)

func (h *Handler) LabReportPage(w http.ResponseWriter, r *http.Request) {
	discID, err := strconv.Atoi(r.PathValue("discID"))
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	labID, err := strconv.Atoi(r.PathValue("labID"))
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	discipline, err := h.db.GetDisciplineByID(discID)
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	lab, err := h.db.GetLabByID(labID)
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	student := models.GetUserFromContext(r.Context())

	report, err := h.db.GetLabReport(student.ID, labID)
	if err != nil {
		report = &models.LabReport{}
	}

	report.Lab = *lab
	report.Discipline = *discipline

	views.LabReportPage(*report).Render(r.Context(), w)
}

func (h *Handler) UploadLabReport(w http.ResponseWriter, r *http.Request) {
	discID, err := strconv.Atoi(r.PathValue("discID"))
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	labID, err := strconv.Atoi(r.PathValue("labID"))
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	discipline, err := h.db.GetDisciplineByID(discID)
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	_, err = h.db.GetLabByID(labID)
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	student := models.GetUserFromContext(r.Context())

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		alerts.FlashError(w, r, "Ошибка чтения формы")
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
		log.Printf("Failed to parse form: %v", err)
		return
	}

	report, err := h.db.GetLabReport(student.ID, labID)
	if err != nil {
		report = &models.LabReport{
			StudentID:    student.ID,
			DisciplineID: discipline.ID,
			LabID:        labID,
			Comment:      r.Form.Get("comment"),
			Status:       "submitted",
		}
		report.UploadedAt = time.Now()
		report.UpdatedAt = time.Now()
	} else {
		report.UpdatedAt = time.Now()
		report.Comment = r.Form.Get("comment")

		log.Printf("equal %v", !report.UpdatedAt.Equal(report.UploadedAt))
	}

	files := r.MultipartForm.File["report_files"]
	storedFiles := []models.StoredFile{}
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		storedFile, err := h.fs.SaveReportFile(
			report,
			file,
			fileHeader,
		)
		if err != nil {
			alerts.FlashError(w, r, "Ошибка сохранения файла")
			log.Printf("Failed to save file: %v", err)
			http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
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

	report.Files = append(report.Files, storedFiles...)

	if report.ID == 0 {
		err = h.db.AddLabReport(report)
		if err != nil {
			alerts.FlashError(w, r, "Ошибка сохранения отчета")
			log.Printf("Failed to add report ot db: %v", err)
			http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
			return
		}
	} else {
		err = h.db.UpdateReport(report)
		if err != nil {
			alerts.FlashError(w, r, "Ошибка обновления отчета")
			log.Printf("Failed to add report ot db: %v", err)
			http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
			return
		}

		err = h.db.AddReportFiles(report.ID, storedFiles)
		if err != nil {
			alerts.FlashError(w, r, "Ошибка обновления отчета")
			log.Printf("Failed to add report ot db: %v", err)
			http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
			return
		}
	}

	alerts.FlashSuccess(w, r, "Отчет отправлен")
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}
