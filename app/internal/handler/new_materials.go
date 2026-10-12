package handler

import (
	"errors"
	"fmt"
	"net/http"
	"scavenger/internal/auth"
	"scavenger/internal/csrf"
	"scavenger/internal/domain"
	"scavenger/internal/repository"
	"scavenger/internal/reqlog"
	"scavenger/internal/service"
	"scavenger/web/templates/pages"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// ---------- Note ----------
func (h *Handler) MaterialNewNotePage(w http.ResponseWriter, r *http.Request) {
	d, u, ok := h.ownedDiscipline(w, r)
	if !ok {
		return
	}
	csrf := csrf.TokenFromCtx(r.Context())
	pages.MaterialNewNotePage(u, *d, csrf, "").Render(r.Context(), w)
}

func (h *Handler) MaterialNewNoteSubmit(w http.ResponseWriter, r *http.Request) {
	log := reqlog.From(r.Context())
	d, u, ok := h.ownedDiscipline(w, r)
	if !ok {
		return
	}
	if err := r.ParseForm(); err != nil {
		log.Warn("parse form", "err", err)
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	_, err := h.svc.Material.CreateNote(r.Context(), service.MaterialNoteInput{
		DisciplineID: d.ID,
		Title:        r.FormValue("title"),
		Description:  r.FormValue("description"),
		Visible:      r.FormValue("visible") != "",
	})
	if err != nil {
		h.renderNewFormError(w, r, u, *d, "note", err)
		return
	}

	Redirect(w, r, fmt.Sprintf("/discipline/%d", d.ID))
}

// ---------- Theory ----------
func (h *Handler) MaterialNewTheoryPage(w http.ResponseWriter, r *http.Request) {
	d, u, ok := h.ownedDiscipline(w, r)
	if !ok {
		return
	}
	csrf := csrf.TokenFromCtx(r.Context())
	pages.MaterialNewTheoryPage(u, *d, csrf, "").Render(r.Context(), w)
}

func (h *Handler) MaterialNewTheorySubmit(w http.ResponseWriter, r *http.Request) {
	log := reqlog.From(r.Context())
	d, u, ok := h.ownedDiscipline(w, r)
	if !ok {
		return
	}
	if err := r.ParseForm(); err != nil {
		log.Warn("parse form", "err", err)
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	_, err := h.svc.Material.CreateTheory(r.Context(), service.MaterialTheoryInput{
		DisciplineID: d.ID,
		Title:        r.FormValue("title"),
		Content:      r.FormValue("content"),
		Visible:      r.FormValue("visible") != "",
	})
	if err != nil {
		h.renderNewFormError(w, r, u, *d, "theory", err)
		return
	}

	Redirect(w, r, fmt.Sprintf("/discipline/%d", d.ID))
}

// ---------- Practical ----------
func (h *Handler) MaterialNewPracticalPage(w http.ResponseWriter, r *http.Request) {
	d, u, ok := h.ownedDiscipline(w, r)
	if !ok {
		return
	}
	csrf := csrf.TokenFromCtx(r.Context())
	pages.MaterialNewPracticalPage(u, *d, csrf, "").Render(r.Context(), w)
}

func (h *Handler) MaterialNewPracticalSubmit(w http.ResponseWriter, r *http.Request) {
	log := reqlog.From(r.Context())
	d, u, ok := h.ownedDiscipline(w, r)
	if !ok {
		return
	}
	if err := r.ParseForm(); err != nil {
		log.Warn("parse form", "err", err)
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	number, _ := strconv.Atoi(r.FormValue("number"))
	deadline, err := time.ParseInLocation("2006-01-02T15:04", r.FormValue("deadline"), time.Local)
	if err != nil {
		log.Warn("parse time", "err", err)
		h.renderNewFormError(w, r, u, *d, "practical", errors.New("некорректный дедлайн"))
		return
	}

	_, err = h.svc.Material.CreatePractical(r.Context(), service.MaterialPracticalInput{
		DisciplineID: d.ID,
		Title:        r.FormValue("title"),
		Content:      r.FormValue("content"),
		Number:       number,
		Deadline:     deadline,
		Visible:      r.FormValue("visible") != "",
	})
	if err != nil {
		h.renderNewFormError(w, r, u, *d, "practical", err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/discipline/%d", d.ID), http.StatusSeeOther)
}

// ---------- shared ----------
func (h *Handler) renderNewFormError(w http.ResponseWriter, r *http.Request, u *domain.User, d domain.Discipline, kind string, err error) {
	log := reqlog.From(r.Context())
	log.Warn("create material", "kind", kind, "err", err)

	msg := "Не удалось создать материал. Проверьте поля и попробуйте снова."
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		msg = "Проверьте заполнение полей."
	case errors.Is(err, repository.ErrNotFound):
		msg = "Дисциплина не найдена."
	}

	csrf := csrf.TokenFromCtx(r.Context())
	switch kind {
	case "note":
		pages.MaterialNewNotePage(u, d, csrf, msg).Render(r.Context(), w)
	case "theory":
		pages.MaterialNewTheoryPage(u, d, csrf, msg).Render(r.Context(), w)
	case "practical":
		pages.MaterialNewPracticalPage(u, d, csrf, msg).Render(r.Context(), w)
	}
}

// internal/handler/material_new.go
func (h *Handler) ownedDiscipline(w http.ResponseWriter, r *http.Request) (*domain.Discipline, *domain.User, bool) {
	log := reqlog.From(r.Context())
	u := auth.UserFromCtx(r.Context())

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return nil, nil, false
	}

	d, err := h.svc.Discipline.ByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.NotFound(w, r)
			return nil, nil, false
		}
		log.Error("discipline by id", "err", err)
		http.Error(w, "internal", http.StatusInternalServerError)
		return nil, nil, false
	}

	if d.TeacherID != u.ID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return nil, nil, false
	}

	return d, u, true
}
