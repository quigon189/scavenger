package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
)

func (h *Handler) GetFile(w http.ResponseWriter, r *http.Request) {
	storedFile, err := h.db.GetStoredFileByURL(r.URL.Path)
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	data, err := h.fs.GetFile(storedFile.Path)
	if err != nil {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	if filepath.Ext(storedFile.Filename) == ".pdf" {
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "inline")
	} else {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", storedFile.Filename))
	}

	w.Write(data)
}
