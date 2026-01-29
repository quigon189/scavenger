package middlewares

import (
	"context"
	"data-service/internal/session"
	"net/http"
)

type AuthMiddleware struct {
	session *session.SessionService
}

func NewAuthMiddleware(session *session.SessionService) *AuthMiddleware {
	return &AuthMiddleware{session: session}
}

func (h *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Header.Get("X-Session-ID")
		if sessionID != "" {
			session, err := h.session.GetSession(r.Context(), sessionID)
			if err == nil {
				ctx := r.Context()
				ctx = context.WithValue(ctx, "session_id", sessionID)
				ctx = context.WithValue(ctx, "session", session)

				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}
}
