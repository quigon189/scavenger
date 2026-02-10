package middlewares

import (
	"context"
	"log"
	"net/http"
)

func (m *Middleware) AdminRequire(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := m.session.GetSession(r)
		if err != nil {
			m.session.DeleteSession(w, r)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			log.Printf("WARN Failed to get session: %v", err)
			return
		}

		ctx := context.WithValue(r.Context(), "session", session)

		if session.User.Role != "admin" {
			m.session.FlashError(w, r, "Доступ запрещен. Требуются права администратора")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))

	}
}
