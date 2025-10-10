package handlers

import (
	"log"
	"net/http"
	"scavenger/internal/alerts"
	"scavenger/views"
	"strconv"
	"strings"

	"github.com/yuin/goldmark"
)

func (h *Handler) LabMarkdownPage(w http.ResponseWriter, r *http.Request) {
	labID, err := strconv.Atoi(r.PathValue("labID"))
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	discID, err := strconv.Atoi(r.PathValue("discID"))
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	var prevURL string

	if r.Referer() != "" {
		prevURL = r.Referer()
	} else {
		prevURL = "/"
	}

	lab, err := h.db.GetLabByID(labID)
	if err != nil {
		alerts.FlashError(w, r, "Работа не найдена")
		log.Printf("Failed to load lab from db: %v", err)
		http.Redirect(w, r, prevURL, http.StatusSeeOther)
		return
	}

	disc, err := h.db.GetDisciplineByID(discID)
	if err != nil {
		alerts.FlashError(w, r, "Дисциплина не найдена")
		log.Printf("Failed to load discipline from db: %v", err)
		http.Redirect(w, r, prevURL, http.StatusSeeOther)
		return
	}

	var htmlContent string

	if lab.MDFile.Path != "" {
		content, err := h.fs.GetFile(lab.MDFile.Path)
		if err != nil {
			alerts.FlashWarning(w, r, "Не удалость прочитать markdown файл")
		} else {
			var buf strings.Builder

			err := goldmark.Convert(content, &buf)
			if err != nil {
				alerts.FlashWarning(w, r, "Не удалось конвертировать md файл")
			} else {
				htmlContent = buf.String()
			}
		}
	}

	views.LabMarkdownPage(*lab, *disc, htmlContent).Render(r.Context(), w)
}
