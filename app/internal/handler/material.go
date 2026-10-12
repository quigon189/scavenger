package handler

import (
	"errors"
	"fmt"
	"net/http"
	"scavenger/internal/auth"
	"scavenger/internal/domain"
	"scavenger/internal/markdown"
	"scavenger/internal/repository"
	"scavenger/internal/reqlog"
	"scavenger/web/templates/components"
	"scavenger/web/templates/pages"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) MaterialView(w http.ResponseWriter, r *http.Request) {
	log := reqlog.From(r.Context())
	u := auth.UserFromCtx(r.Context())
	id := chi.URLParam(r, "id")

	m, err := h.svc.Material.ByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		log.Error("material by id", "err", err)
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}

	if u.Role == domain.RoleStudent && !m.GetVisible() {
		http.NotFound(w, r)
		return
	}

	if m.Kind() == domain.MaterialNoteType {
		Redirect(w, r, fmt.Sprintf("/discipline/%d", m.GetDisciplineID()))
		return
	}

	disc, err := h.svc.Discipline.ByID(r.Context(), m.GetDisciplineID())
	if err != nil {
		log.Error("discipline by id", "err", err)
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}

	group, _ := h.svc.Group.ByID(r.Context(), disc.GroupID)
	teacher, _ := h.svc.User.ByID(r.Context(), disc.TeacherID)

	view := components.MaterialPageView{
		Discipline: *disc,
		Group:      group,
		Teacher:    teacher,
		User:       u,
	}

	switch u.Role {
	case domain.RoleStudent:
		view.Prev, view.Next, err = h.svc.Material.GetNieghbours(r.Context(), m, false)
	case domain.RoleTeacher:
		view.Prev, view.Next, err = h.svc.Material.GetNieghbours(r.Context(), m, true)
	}
	if err != nil {
		log.Error("material nieghbours", "err", err)
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}

	switch v := m.(type) {
	case *domain.TheoryMaterial:
		html, err := markdown.Render(v.Content)
		if err != nil {
			log.Error("markdown", "err", err)
			http.Error(w, "internal", http.StatusInternalServerError)
			return
		}
		view.Theory = &components.TheoryPageBody{Material: *v, HTML: html}
	case *domain.PracticalMaterial:
		html, err := markdown.Render(v.Content)
		if err != nil {
			log.Error("markdown", "err", err)
			http.Error(w, "internal", http.StatusInternalServerError)
			return
		}
		grade := 5
		progress := components.Progress{
			Total: 10,
			Done:  7,
		}
		view.Practical = &components.PracticalPageBody{
			Material: *v,
			HTML:     html,
			Grade:    &grade,
			Progress: &progress,
		}
	}

	pages.MaterialPage(u, view).Render(r.Context(), w)
}

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
