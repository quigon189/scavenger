package handler

import (
	"log/slog"
	"net/http"
	"scavenger/internal/auth"
	"scavenger/internal/domain"
	"scavenger/internal/reqlog"
	"slices"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func RequestLogger(base *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()

			log := base.With(
				slog.String("request_id", middleware.GetReqID(r.Context())),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
			)

			ctx := reqlog.WithLogger(r.Context(), log)
			r = r.WithContext(ctx)

			defer func() {
				lvl := slog.LevelInfo
				switch {
				case ww.Status() >= 500:
					lvl = slog.LevelError
				case ww.Status() >= 400:
					lvl = slog.LevelWarn
				}

				log.LogAttrs(ctx, lvl, "http request",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.Duration("duration", time.Since(start)),
				)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := reqlog.From(r.Context())
		c, err := r.Cookie(auth.SessionCookie)
		if err != nil || c.Value == "" {
			next.ServeHTTP(w,r)
			return
		}
		u, err := h.svc.Auth.UserBySession(r.Context(), c.Value)
		if err != nil {
			log.Warn("user session", "err", err)
			auth.ClearSessionCookie(w)
			next.ServeHTTP(w, r)
			return
		}
		ctx := auth.WithUser(r.Context(), u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if auth.UserFromCtx(r.Context()) == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RequireRole (roles ...domain.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u := auth.UserFromCtx(r.Context())
			if u == nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			if !slices.Contains(roles, u.Role) {
				http.Error(w, "forbidden", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
