package middlewares

import (
	"context"
	"log"
	"net/http"
	"web-service/internal/session"
)

type AuthMiddleware struct {
	session *session.SessionService
}

func NewAuthMiddleware(session *session.SessionService) *AuthMiddleware {
	return &AuthMiddleware{
		session: session,
	}
}

func (m *AuthMiddleware) SessionRequire(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, sessionID, err := m.session.GetSession(r)
		if err != nil {
			m.session.DeleteSession(w)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			log.Printf("WARN Failed to get session %s: %v", sessionID, err)
			return
		}

		ctx := context.WithValue(r.Context(), "session_id", sessionID)
		ctx = context.WithValue(ctx, "session", session)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
