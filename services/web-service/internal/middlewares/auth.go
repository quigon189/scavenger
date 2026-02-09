package middlewares

import (
	"context"
	"log"
	"net/http"

	"web-service/internal/session"
)

type Middleware struct {
	session *session.SessionService
}

func NewMiddleware(session *session.SessionService) *Middleware {
	return &Middleware{
		session: session,
	}
}

func (m *Middleware) SessionRequire(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, sessionID, err := m.session.GetSession(r)
		if err != nil {
			m.session.DeleteSession(w, r)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			log.Printf("WARN Failed to get session %s: %v", sessionID, err)
			return
		}

		ctx := context.WithValue(r.Context(), "session_id", sessionID)
		ctx = context.WithValue(ctx, "session", session)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (m *Middleware) ActiveRequire(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, sessionID, err := m.session.GetSession(r)
		if err != nil {
			m.session.DeleteSession(w, r)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			log.Printf("WARN Failed to get session %s: %v", sessionID, err)
			return
		}

		ctx := context.WithValue(r.Context(), "session_id", sessionID)
		ctx = context.WithValue(ctx, "session", session)

		if session.User.Status == "pending" {
			http.Redirect(w, r, "/pending", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))

	}
}

func (m *Middleware) FlashHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flashes := m.session.GetFlashes(w, r)
		ctx := context.WithValue(r.Context(), "flashes", flashes)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
