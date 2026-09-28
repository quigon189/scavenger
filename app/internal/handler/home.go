package handler

import (
	"net/http"
	"scavenger/internal/auth"
	"scavenger/web/templates/pages"
)

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
    u := auth.UserFromCtx(r.Context())
    pages.Home(u).Render(r.Context(), w)
}
