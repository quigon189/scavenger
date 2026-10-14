package handler

import (
	"errors"
	"fmt"
	"net/http"
	"scavenger/internal/auth"
	"scavenger/internal/csrf"
	"scavenger/internal/domain"
	"scavenger/internal/markdown"
	"scavenger/internal/repository"
	"scavenger/internal/reqlog"
	"scavenger/internal/service"
	"scavenger/web/templates/components"
	"scavenger/web/templates/pages"
	"scavenger/web/templates/viewmodels"
	"strconv"
	"time"

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

	view := viewmodels.MaterialPageView{
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
		view.Theory = &viewmodels.TheoryPageBody{Material: *v, HTML: html}
	case *domain.PracticalMaterial:
		html, err := markdown.Render(v.Content)
		if err != nil {
			log.Error("markdown", "err", err)
			http.Error(w, "internal", http.StatusInternalServerError)
			return
		}
		grade := 5
		progress := viewmodels.Progress{
			Total: 10,
			Done:  7,
		}
		view.Practical = &viewmodels.PracticalPageBody{
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

func (h *Handler) MaterialFormNew(w http.ResponseWriter, r *http.Request) {
	u := auth.UserFromCtx(r.Context())
	log := reqlog.From(r.Context())

	typ, ok := parseMaterialType(chi.URLParam(r, "type"))
	if !ok {
		http.NotFound(w, r)
		return
	}

	disciplineID, err := strconv.ParseInt(r.URL.Query().Get("discipline_id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid discipline_id", http.StatusBadRequest)
		return
	}

	discipline, err := h.svc.Discipline.ByID(r.Context(), disciplineID)
	if err != nil {
		log.Error("failed to get discipline", "discipline_id", disciplineID, "err", err)
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	if discipline.TeacherID != u.ID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	form := viewmodels.MaterialFormData{
		DisciplineID: discipline.ID,
		Type:         string(typ),
		Visible:      true,
		CSRF: csrf.TokenFromCtx(r.Context()),
	}

	pages.MaterialForm(u, *discipline, form).Render(r.Context(), w)
}

func (h *Handler) MaterialFormNewSubmit(w http.ResponseWriter, r *http.Request) {
	log := reqlog.From(r.Context())
	u := auth.UserFromCtx(r.Context())

	typ, ok := parseMaterialType(chi.URLParam(r, "type"))
	if !ok {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		log.Warn("parse form", "err", err)
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	var form viewmodels.MaterialFormData
	form.DisciplineID, _ = strconv.ParseInt(r.FormValue("discipline_id"), 10, 64)
	form.Type = string(typ)
	form.Title = r.FormValue("title")
	form.Visible = r.FormValue("visible") != ""
	form.Description = r.FormValue("description")
	form.Content = r.FormValue("content")
	form.Deadline = r.FormValue("deadline")
	form.Number, _ = strconv.Atoi(r.FormValue("number"))

	attachmentFileIDs := r.MultipartForm.Value["attachment_file_ids"]

	discipline, err := h.svc.Discipline.ByID(r.Context(), form.DisciplineID) 
	if err != nil {
		log.Error("get discipline", "id", form.DisciplineID, "err", err)
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}

	if discipline.TeacherID != u.ID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	files, err := h.svc.File.ByIDs(r.Context(), attachmentFileIDs)
	if err != nil {
		log.Error("get files", "ids", attachmentFileIDs, "err", err)
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}

	for _, f := range files {
		form.Attachments = append(form.Attachments, *f)
	}

	var createErr error
	switch typ {
	case domain.MaterialNoteType:
		_, createErr = h.svc.Material.CreateNote(r.Context(), service.MaterialNoteInput{
			DisciplineID: discipline.ID,
			Title: form.Title,
			Description: form.Description,
			Visible: form.Visible,
		})
	case domain.MaterialTheoryType:
		_, createErr = h.svc.Material.CreateTheory(r.Context(), service.MaterialTheoryInput{
			DisciplineID: discipline.ID,
			Title: form.Title,
			Content: form.Content,
			Visible: form.Visible,
			AttachmentFileIDs: attachmentFileIDs,	
		})
	case domain.MaterialPracticalType:
		deadline, perr := time.ParseInLocation("2006-01-02T15:04", form.Deadline, time.Local)
		if perr != nil {
			createErr = errors.New("некорректный дедлайн")
			break
		}
		_, createErr = h.svc.Material.CreatePractical(r.Context(), service.MaterialPracticalInput{
			DisciplineID: discipline.ID,
			Title: form.Title,
			Content: form.Content,
			Visible: form.Visible,
			Number: form.Number,
			Deadline: deadline,
			AttachmentFileIDs: attachmentFileIDs,
		})
	}

	if createErr != nil {
		form.ErrMsg = humanizeMaterialError(createErr)
		form.CSRF = csrf.TokenFromCtx(r.Context())
		pages.MaterialForm(u, *discipline, form)
		return
	}
	Redirect(w, r, fmt.Sprintf("/discipline/%d", discipline.ID))
}

func (h *Handler) MaterialAttachmentUpload(w http.ResponseWriter, r *http.Request) {
	log := reqlog.From(r.Context())
	u := auth.UserFromCtx(r.Context())

	if err := r.ParseMultipartForm(h.cfg.MaxUploadSize); err != nil {
		log.Warn("parse multipart", "err", err)
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	fhs := r.MultipartForm.File["attachment_files"]
	if len(fhs) == 0 {
		http.Error(w, "no files", http.StatusBadRequest)
		return
	}

	for _, fh := range fhs {
		src, err := fh.Open()
		if err != nil {
			log.Error("open multipart file", "err", err)
			continue
		}
		f, err := h.svc.File.Upload(r.Context(), service.FileUploadInput{
			OwnerID:      u.ID,
			OriginalName: fh.Filename,
			ContentType:  fh.Header.Get("Content-Type"),
			Size:         fh.Size,
			Reader:       src,
		})
		src.Close()
		if err != nil {
			log.Error("upload file", "err", err)
			continue
		}
		components.AttachmentItem(viewmodels.AttachmentItemView{File: *f}).
			Render(r.Context(), w)
	}
}

// POST /materials/attachment/{id}/delete
func (h *Handler) MaterialAttachmentDelete(w http.ResponseWriter, r *http.Request) {
	log := reqlog.From(r.Context())
	u := auth.UserFromCtx(r.Context())
	id := chi.URLParam(r, "id")

	f, err := h.svc.File.ByID(r.Context(), id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if f.OwnerID != u.ID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	if err := h.svc.File.Delete(r.Context(), id); err != nil {
		log.Error("delete file", "err", err)
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	// Пустой ответ — htmx удалит элемент через outerHTML-swap.
	w.WriteHeader(http.StatusNoContent)
}

func humanizeMaterialError(err error) string {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		return "Проверьте заполнение полей."
	case errors.Is(err, repository.ErrNotFound):
		return "Связанная запись не найдена."
	}
	return err.Error()
}
