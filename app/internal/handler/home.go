package handler

import (
	"cmp"
	"fmt"
	"net/http"
	"scavenger/internal/auth"
	"scavenger/internal/domain"
	"scavenger/internal/reqlog"
	"scavenger/web/templates/components"
	"scavenger/web/templates/pages"
	"slices"
)

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	log := reqlog.From(r.Context())
	u := auth.UserFromCtx(r.Context())
	if u.Role == domain.RoleStudent {
		group, err := h.svc.Group.ByStudentID(r.Context(), u.ID)
		if err != nil {
			log.Error("failed to get student group", "err", err)
			http.Error(w, "internal", http.StatusInternalServerError)
			return
		}
		disciplines, err := h.svc.Discipline.ByGroupID(r.Context(), group.ID)
		if err != nil {
			log.Error("failed to get group disciplines", "group_id", group.ID, "err", err)
		}
		slices.SortFunc(disciplines, func(a, b domain.Discipline) int {
			return cmp.Compare(a.Title, b.Title)
		})
		items := make([]components.DisciplineCardView, 0, len(disciplines))
		for _, d := range disciplines {
			teacher, _ := h.svc.User.ByID(r.Context(), d.TeacherID)
			items = append(items, components.DisciplineCardView{
				Discipline: d,
				Teacher: teacher,
			})
		}
		pages.StudentHome(u, group, items).Render(r.Context(), w)
		return
	}
	if u.Role == domain.RoleTeacher {
		disciplines, err := h.svc.Discipline.ByTeacherID(r.Context(), u.ID)
		if err != nil {
			log.Error("failed to get teacher disciplines", "teacher_id", u.ID, "err", err)
			http.Error(w, "internal", http.StatusInternalServerError)
			return
		}
		groupIDs := uniqueInt64(disciplines, func(d domain.Discipline) int64 {
			return d.GroupID
		})
		var groups []domain.Group
		for _, gid := range groupIDs {
			group, err := h.svc.Group.ByID(r.Context(), gid)
			if err != nil {
				log.Warn("failed to get group", "group_id", gid)
				continue
			}
			groups = append(groups, *group)
		}

		slices.SortFunc(groups, func(a, b domain.Group) int {
			return cmp.Compare(a.Number, b.Number)
		})

		byGroup := map[int64][]domain.Discipline{}
		for _, d := range disciplines {
			byGroup[d.GroupID] = append(byGroup[d.GroupID], d)
		}

		sections := make([]components.DisciplineGroupView, 0, len(groups))
		for _, g := range groups {
			items := make([]components.DisciplineCardView, 0, len(byGroup[g.ID]))
			for _, d := range byGroup[g.ID] {
				items = append(items, components.DisciplineCardView{
					Discipline: d,
					Group:      &g,
				})
			}
			sections = append(sections, components.DisciplineGroupView{
				ID:       fmt.Sprintf("group-%d", g.Number),
				Title:    g.Name(),
				Subtitle: plural(len(items), "дисциплина", "дисциплины", "дисциплин"),
				Items:    items,
			})
		}
		pages.TeacherHome(u, sections).Render(r.Context(), w)
		return
	}

	pages.Home(u).Render(r.Context(), w)
}
