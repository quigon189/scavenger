package handlers

import (
	"context"
	"log"
	"net/http"
	"scavenger/internal/alerts"
)

func (h *Handler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !h.authService.IsAuthenticated(r) {
			h.Login(w,r)
			return
		}
		user := h.authService.GetUser(r)

		ctx := context.WithValue(r.Context(), "user", user)

		log.Printf("Login user: %+v", user)

		next.ServeHTTP(w,r.WithContext(ctx))
	}
}

func (h *Handler) AdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return h.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		role := h.authService.GetUserRole(r)
		if role != "admin" {
			http.Redirect(w,r,"/", http.StatusSeeOther)
		}
		next.ServeHTTP(w,r)
	})
}

func (h *Handler) StudentMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return h.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		role := h.authService.GetUserRole(r)
		if role != "student" {
			http.Redirect(w,r,"/", http.StatusSeeOther)
		}
		next.ServeHTTP(w,r)
	})
}

func (h *Handler) AlertMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		als := alerts.GetAlerts(w,r)

		ctx := context.WithValue(r.Context(), alerts.AlertKey, als)

		next.ServeHTTP(w,r.WithContext(ctx))
	})
}
