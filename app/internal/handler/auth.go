package handler

import (
	"errors"
	"net/http"
	"scavenger/internal/auth"
	"scavenger/internal/reqlog"
	"scavenger/internal/service"
	"scavenger/web/templates/pages"
	"time"
)

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	pages.LoginPage("").Render(r.Context(), w)
}

func (h *Handler) LoginSubmit(w http.ResponseWriter, r *http.Request) {
	log := reqlog.From(r.Context())
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")

	sid, _, err := h.svc.Auth.Login(r.Context(), service.LoginInput{
		Email: email, Password: password,
	})
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredetials) {
			pages.LoginPage("Неверный email или пароль").Render(r.Context(), w)
			return
		}
		if errors.Is(err, service.ErrUnactiveUser) {
			pages.LoginPage("Пользователь неактивен, обратитесь к администратору"). Render(r.Context(), w)
			return
		}
		log.Error("login", "err", err)
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}

	auth.SetSessionCookie(w, sid, time.Now().Add(1*time.Hour), false)

	Redirect(w, r, "/")
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
    if c, err := r.Cookie(auth.SessionCookie); err == nil {
        _ = h.svc.Auth.Logout(r.Context(), c.Value)
    }
    auth.ClearSessionCookie(w)
    http.Redirect(w, r, "/login", http.StatusSeeOther)
}
