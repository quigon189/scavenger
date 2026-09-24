package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"
)

const SessionCookie = "sid"

func NewSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func SetSessionCookie(w http.ResponseWriter, id string, expires time.Time, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name: SessionCookie,
		Value: id,
		Path: "/",
		Expires: expires,
		HttpOnly: true,
		Secure: secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name: SessionCookie,
		Value: "",
		Path: "/",
		MaxAge: -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
