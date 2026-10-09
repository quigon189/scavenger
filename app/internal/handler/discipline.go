package handler

import (
	"cmp"
	"net/http"
	"scavenger/internal/auth"
	"scavenger/internal/domain"
	"scavenger/internal/reqlog"
	"scavenger/web/templates/components"
	"scavenger/web/templates/pages"
	"slices"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) Discipline(w http.ResponseWriter, r *http.Request) {
	disciplineIDStr := chi.URLParam(r, "id")

	disciplineID, err := strconv.ParseInt(disciplineIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	u := auth.UserFromCtx(r.Context())
	switch u.Role {
	case domain.RoleStudent:
		h.StudentDiscipline(w, r, disciplineID)
	case domain.RoleTeacher:
		h.TeacherDiscipline(w, r, disciplineID)
	default:
		Redirect(w, r, "/")
	}
}

func (h *Handler) TeacherDiscipline(w http.ResponseWriter, r *http.Request, disciplineID int64) {
	log := reqlog.From(r.Context())
	u := auth.UserFromCtx(r.Context())
	discipline, err := h.svc.Discipline.ByID(r.Context(), disciplineID)
	if err != nil {
		log.Error("failed to get discipline", "id", disciplineID, "err", err)
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	if discipline.TeacherID != u.ID {
		log.Warn("invalid teacher id for discipline", "user", u, "discipline", discipline)
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	materials, err := h.svc.Material.ListByDisciplineID(r.Context(), disciplineID)
	if err != nil {
		log.Error("failed to get dicipline materials", "discipline_id", discipline.ID, "err", err)
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	group, _ := h.svc.Group.ByID(r.Context(), discipline.GroupID)

	dcv := components.DisciplineCardView{
		Discipline: *discipline,
		Group:      group,
	}

	slices.SortFunc(materials, func(a, b domain.Material) int {
		return cmp.Compare(a.GetDisplayOrder(), b.GetDisplayOrder())
	})

	materialCardsViews := make([]components.MaterialCardView, 0, len(materials))
	for _, m := range materials {
		switch v := m.(type) {
		case domain.NoteMaterial:
			materialCardsViews = append(materialCardsViews, components.MaterialCardView{
				Note: &components.NoteCardView{Material: v},
			})
		case domain.TheoryMaterial:
			materialCardsViews = append(materialCardsViews, components.MaterialCardView{
				Theory: &components.TheoryCardView{Material: v},
			})
		case domain.PracticalMaterial:
			materialCardsViews = append(materialCardsViews, components.MaterialCardView{
				Practical: &components.PracticalCardView{
					Material: v,
					Progress: &components.Progress{
						Done:  7,
						Total: 10,
					},
				},
			})
		}
	}

	pages.DisciplinePage(u, dcv, materialCardsViews).Render(r.Context(), w)
}

func (h *Handler) StudentDiscipline(w http.ResponseWriter, r *http.Request, disciplineID int64) {
	log := reqlog.From(r.Context())
	u := auth.UserFromCtx(r.Context())
	discipline, err := h.svc.Discipline.ByID(r.Context(), disciplineID)
	if err != nil {
		log.Error("failed to get discipline", "id", disciplineID, "err", err)
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	group, err := h.svc.Group.ByStudentID(r.Context(), u.ID)
	if err != nil {
		log.Error("failed to get group", "id", u.ID, "err", err)
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	if discipline.GroupID != group.ID {
		log.Warn("invalid group id for discipline", "group", group, "discipline", discipline)
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	materials, err := h.svc.Material.ListByDisciplineID(r.Context(), disciplineID)
	if err != nil {
		log.Error("failed to get dicipline materials", "discipline_id", discipline.ID, "err", err)
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	teacher, _ := h.svc.User.ByID(r.Context(), discipline.TeacherID)

	dcv := components.DisciplineCardView{
		Discipline: *discipline,
		Teacher:    teacher,
	}

	slices.SortFunc(materials, func(a, b domain.Material) int {
		return cmp.Compare(a.GetDisplayOrder(), b.GetDisplayOrder())
	})

	materialCardsViews := make([]components.MaterialCardView, 0, len(materials))
	for _, m := range materials {
		switch v := m.(type) {
		case domain.NoteMaterial:
			materialCardsViews = append(materialCardsViews, components.MaterialCardView{
				Note: &components.NoteCardView{Material: v},
			})
		case domain.TheoryMaterial:
			materialCardsViews = append(materialCardsViews, components.MaterialCardView{
				Theory: &components.TheoryCardView{Material: v},
			})
		case domain.PracticalMaterial:
			grade := 5
			materialCardsViews = append(materialCardsViews, components.MaterialCardView{
				Practical: &components.PracticalCardView{
					Material: v,
					Grade: &grade,
				},
			})
		}
	}

	pages.DisciplinePage(u, dcv, materialCardsViews).Render(r.Context(), w)

}
