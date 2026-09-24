package csrf

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"scavenger/internal/reqlog"
)

const CookieName = "csrf"
const HeaderName = "X-CSRF-Token"
const FormField = "_csrf"

type ctxKey int

const tokenKey ctxKey = 1

func TokenFromCtx(ctx context.Context) string {
	s, _ := ctx.Value(tokenKey).(string)
	return s
}

func WithToken(ctx context.Context, t string) context.Context {
	return context.WithValue(ctx, tokenKey, t)
}

func NewToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string
		log := reqlog.From(r.Context())
		if c, err := r.Cookie(CookieName); err == nil && c.Value != "" {
			token = c.Value
		} else {
			t, err := NewToken()
			if err != nil {
				http.Error(w, "csrf token", http.StatusInternalServerError)
				log.Error("csrf token", "err", err)
				return
			}
			token = t
			http.SetCookie(w, &http.Cookie{
				Name:     CookieName,
				Value:    token,
				Path:     "/",
				HttpOnly: false,
				SameSite: http.SameSiteLaxMode,
			})
		}
		next.ServeHTTP(w, r.WithContext(WithToken(r.Context(), token)))
	})
}

func Verify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := reqlog.From(r.Context())
		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			next.ServeHTTP(w, r)
			return
		}
		c, err := r.Cookie(CookieName)
		if err != nil || c.Value == "" {
			http.Error(w, "csrf token missing", http.StatusForbidden)
			if err != nil {
				log.Warn("csrf token missing", "err", err)
			} else {
				log.Warn("csrf token missing")
			}
			return
		}
		got := r.Header.Get(HeaderName)
		if got == "" {
			got = r.FormValue(FormField)
		}
		if got != c.Value {
			http.Error(w, "csrf mismath", http.StatusForbidden)
			log.Warn("csrf mismatch")
			return
		}
		next.ServeHTTP(w, r)
	})
}
