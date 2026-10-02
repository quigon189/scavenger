package handler

import (
	"net/http"
	"scavenger/internal/auth"
	"scavenger/internal/domain"
	"scavenger/internal/reqlog"
	"scavenger/web/templates/pages"
)

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	log := reqlog.From(r.Context())
    u := auth.UserFromCtx(r.Context())
	if u.Role == domain.RoleStudent {
		group, err := h.svc.Group.ByStudentID(r.Context(), u.ID)
		if err != nil {
			log.Warn("student group", "err", err)
		} else {
    		pages.Home(u, group).Render(r.Context(), w)
			return
		}
	}
    pages.Home(u, nil).Render(r.Context(), w)
}
