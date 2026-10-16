package handler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"scavenger/internal/auth"
	"scavenger/internal/csrf"
	"scavenger/internal/domain"
	"scavenger/internal/markdown"
	"scavenger/internal/repository"
	"scavenger/internal/reqlog"
	"scavenger/internal/routes"
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

	var attachmentFileIDs []string
	switch v := m.(type) {
	case *domain.TheoryMaterial:
		html, err := markdown.Render(v.Content)
		if err != nil {
			log.Error("markdown", "err", err)
			http.Error(w, "internal", http.StatusInternalServerError)
			return
		}
		view.Theory = &viewmodels.TheoryPageBody{Material: *v, HTML: html}
		attachmentFileIDs = v.AttachmetFileIDs
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
		attachmentFileIDs = v.AttachmetFileIDs
	}

	if len(attachmentFileIDs) > 0 {
		files, err := h.svc.File.ByIDs(r.Context(), attachmentFileIDs)
		if err != nil {
			log.Error("material attachments", "err", err)
		} else {
			for _, file := range files {
				view.Attachments = append(view.Attachments, *file)
			}
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

func (h *Handler) MaterialShow(w http.ResponseWriter, r *http.Request) {
	h.materialVisible(w, r, true)
}

func (h *Handler) MaterialHide(w http.ResponseWriter, r *http.Request) {
	h.materialVisible(w, r, false)
}

func (h *Handler) materialVisible(w http.ResponseWriter, r *http.Request, visible bool) {
	u := auth.UserFromCtx(r.Context())
	log := reqlog.From(r.Context())

	id := chi.URLParam(r, "id")

	material, err := h.svc.Material.ByID(r.Context(), id)
	if err != nil {
		log.Error("get material", "id", id, "err", err)
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	discipline, err := h.svc.Discipline.ByID(r.Context(), material.GetDisciplineID())
	if err != nil {
		log.Error("get discipline", "id", material.GetDisciplineID(), "err", err)
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	if discipline.TeacherID != u.ID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	material.SetVisible(visible)
	if err := h.svc.Material.Update(r.Context(), material); err != nil {
		log.Error("update discipline", "id", material.GetDisciplineID(), "err", err)
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Refresh", "true")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func (h *Handler) MaterialDelete(w http.ResponseWriter, r *http.Request) {
	u := auth.UserFromCtx(r.Context())
	log := reqlog.From(r.Context())
	id := chi.URLParam(r, "id")

	material, err := h.svc.Material.ByID(r.Context(), id)
	if err != nil {
		log.Error("get material", "id", id, "err", err)
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	discipline, err := h.svc.Discipline.ByID(r.Context(), material.GetDisciplineID())
	if err != nil {
		log.Error("get discipline", "id", material.GetDisciplineID(), "err", err)
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	if discipline.TeacherID != u.ID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if err := h.svc.Material.Delete(r.Context(), material.GetID()); err != nil {
		log.Error("delete discipline", "id", material.GetDisciplineID(), "err", err)
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}

	Redirect(w, r, routes.DisciplineID.Link(discipline.ID))
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
		CSRF:         csrf.TokenFromCtx(r.Context()),
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
	form.Content = r.FormValue("markdown-content")
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
			Title:        form.Title,
			Description:  form.Description,
			Visible:      form.Visible,
		})
	case domain.MaterialTheoryType:
		_, createErr = h.svc.Material.CreateTheory(r.Context(), service.MaterialTheoryInput{
			DisciplineID:      discipline.ID,
			Title:             form.Title,
			Content:           form.Content,
			Visible:           form.Visible,
			AttachmentFileIDs: attachmentFileIDs,
		})
	case domain.MaterialPracticalType:
		deadline, perr := time.ParseInLocation("2006-01-02T15:04", form.Deadline, time.Local)
		if perr != nil {
			createErr = errors.New("некорректный дедлайн")
			break
		}
		_, createErr = h.svc.Material.CreatePractical(r.Context(), service.MaterialPracticalInput{
			DisciplineID:      discipline.ID,
			Title:             form.Title,
			Content:           form.Content,
			Visible:           form.Visible,
			Number:            form.Number,
			Deadline:          deadline,
			AttachmentFileIDs: attachmentFileIDs,
		})
	}

	if createErr != nil {
		form.ErrMsg = humanizeMaterialError(createErr)
		form.CSRF = csrf.TokenFromCtx(r.Context())
		pages.MaterialForm(u, *discipline, form).Render(r.Context(), w)
		return
	}
	Redirect(w, r, fmt.Sprintf("/discipline/%d", discipline.ID))
}

func (h *Handler) MaterialFormEdit(w http.ResponseWriter, r *http.Request) {
	u := auth.UserFromCtx(r.Context())
	log := reqlog.From(r.Context())
	id := chi.URLParam(r, "id")

	material, err := h.svc.Material.ByID(r.Context(), id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	discipline, err := h.svc.Discipline.ByID(r.Context(), material.GetDisciplineID())
	if err != nil {
		log.Error("get discipline", "id", material.GetDisciplineID(), "err", err)
		http.Error(w, "inernal", http.StatusInternalServerError)
		return
	}
	if discipline.TeacherID != u.ID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	form := viewmodels.MaterialFormData{
		ID:           material.GetID(),
		DisciplineID: discipline.ID,
		Type:         string(material.Kind()),
		IsEdit:       true,
		Title:        material.GetTitle(),
		Visible:      material.GetVisible(),
		CSRF:         csrf.TokenFromCtx(r.Context()),
	}

	attachmentFileIDs := []string{}

	switch v := material.(type) {
	case *domain.NoteMaterial:
		form.Description = v.Description
	case *domain.TheoryMaterial:
		form.Content = v.Content
		attachmentFileIDs = v.AttachmetFileIDs
	case *domain.PracticalMaterial:
		form.Content = v.Content
		form.Number = v.Number
		form.Deadline = v.Deadline.Format("2006-01-02T15:04")
		attachmentFileIDs = v.AttachmetFileIDs
	}

	files, err := h.svc.File.ByIDs(r.Context(), attachmentFileIDs)
	if err != nil {
		log.Error("get attachments", "err", err)
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}

	for _, file := range files {
		form.Attachments = append(form.Attachments, *file)
	}

	pages.MaterialForm(u, *discipline, form).Render(r.Context(), w)
}

func (h *Handler) MaterialFormEditSubmit(w http.ResponseWriter, r *http.Request) {
	u := auth.UserFromCtx(r.Context())
	log := reqlog.From(r.Context())
	id := chi.URLParam(r, "id")

	existing, err := h.svc.Material.ByID(r.Context(), id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	discipline, err := h.svc.Discipline.ByID(r.Context(), existing.GetDisciplineID())
	if err != nil {
		log.Error("get discipline", "id", existing.GetDisciplineID(), "err", err)
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	if discipline.TeacherID != u.ID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		log.Warn("parse form", "err", err)
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	var form viewmodels.MaterialFormData
	form.DisciplineID, _ = strconv.ParseInt(r.FormValue("discipline_id"), 10, 64)
	form.Type = string(existing.Kind())
	form.Title = r.FormValue("title")
	form.Visible = r.FormValue("visible") != ""
	form.Description = r.FormValue("description")
	form.Content = r.FormValue("markdown-content")
	form.Deadline = r.FormValue("deadline")
	form.Number, _ = strconv.Atoi(r.FormValue("number"))
	form.CSRF = csrf.TokenFromCtx(r.Context())
	form.IsEdit = true

	attachmentFileIDs := r.MultipartForm.Value["attachment_file_ids"]

	files, err := h.svc.File.ByIDs(r.Context(), attachmentFileIDs)
	if err != nil {
		log.Error("get files", "ids", attachmentFileIDs, "err", err)
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}

	for _, f := range files {
		form.Attachments = append(form.Attachments, *f)
	}

	var updErr error
	switch v := existing.(type) {
	case *domain.NoteMaterial:
		v.Title = form.Title
		v.Description = form.Description
		v.Visible = form.Visible
		updErr = h.svc.Material.Update(r.Context(), v)
	case *domain.TheoryMaterial:
		v.Title = form.Title
		v.Content = form.Content
		v.Visible = form.Visible
		v.AttachmetFileIDs = attachmentFileIDs
		updErr = h.svc.Material.Update(r.Context(), v)
	case *domain.PracticalMaterial:
		deadline, perr := time.ParseInLocation("2006-01-02T15:04", form.Deadline, time.Local)
		if perr != nil {
			updErr = errors.New("некорректный дедлайн")
			break
		}
		v.Title = form.Title
		v.Content = form.Content
		v.Number = form.Number
		v.Deadline = deadline
		v.Visible = form.Visible
		v.AttachmetFileIDs = attachmentFileIDs
		updErr = h.svc.Material.Update(r.Context(), v)
	}

	if updErr != nil {
		form.ErrMsg = humanizeMaterialError(updErr)
		pages.MaterialForm(u, *discipline, form).Render(r.Context(), w)
		return
	}

	Redirect(w, r, routes.DisciplineID.Link(discipline.ID))
}

func (h *Handler) MaterialMarkdownUpload(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		http.NotFound(w, r)
	}
	log := reqlog.From(r.Context())

	if err := r.ParseMultipartForm(h.cfg.MaxUploadSize); err != nil {
		log.Warn("parse multipart", "err", err)
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("mdfile")	
	if err != nil {
		log.Warn("bad mdfile", "err", err)
		http.Error(w, "bad file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	var mdContent bytes.Buffer
	_, err = io.Copy(&mdContent, file)
	if err != nil {
		log.Error("read mdfile", "err", err)
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}

	components.MaterialMarkdownTextarea(mdContent.String()).Render(r.Context(), w)
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

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	components.AttachmentItem(viewmodels.AttachmentItemView{
		File:    *f,
		Deleted: true,
	}).Render(r.Context(), w)
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
