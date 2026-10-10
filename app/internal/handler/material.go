package handler

import (
	"net/http"
	"scavenger/internal/reqlog"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) MaterialMove(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	dir := r.URL.Query().Get("dir")

	shift := 0
	switch dir {
	case "up":
		shift = -1
	case "down":
		shift = 1
	default:
		http.Error(w, "invalid dir", http.StatusBadRequest)
		return
	}

	if err := h.svc.Material.Move(r.Context(), id, shift); err != nil {
		reqlog.From(r.Context()).Error("material move", "err", err)
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}

	// Простейший вариант: перезагрузить страницу на клиенте.
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Refresh", "true")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}
