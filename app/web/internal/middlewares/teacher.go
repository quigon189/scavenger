package middlewares

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func (m *Middleware) TeacherRequire(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := m.session.GetSession(r)
		if err != nil {
			m.session.DeleteSession(w, r)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			log.Printf("WARN Failed to get session: %v", err)
			return
		}

		ctx := context.WithValue(r.Context(), "session", session)

		if session.User.Role != "teacher" {
			m.session.FlashError(w, r, "Доступ запрещен. Требуются права преподавателя", fmt.Errorf("access denied"))
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
