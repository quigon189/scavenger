package middlewares

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
	"web-service/internal/models"
)

type AlertsMiddleware struct {
}

func NewAlertsMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (m *AuthMiddleware) ProcessAlerts(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alerts := getAlerts(r)

		ctx := context.WithValue(r.Context(), "alerts", alerts)
		next.ServeHTTP(w, r.WithContext(ctx))

		setAlerts(w, alerts)
	}
}

func getAlerts(r *http.Request) []models.Alert {
	var alerts []models.Alert
	if r == nil {
		return alerts
	}

	cookie, err := r.Cookie("alerts")
	if err != nil {
		return alerts
	}

	decoded, _ := url.QueryUnescape(cookie.Value)

	json.Unmarshal([]byte(decoded), &alerts)

	return alerts
}

func setAlerts(w http.ResponseWriter, alerts []models.Alert) {
	if len(alerts) == 0 {
		http.SetCookie(w, &http.Cookie{
			Name:     "alerts",
			Value:    "",
			Path:     "/",
			Expires:  time.Now().Add(-1 * time.Minute),
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
		return
	}

	jsonData, _ := json.Marshal(alerts)
	encoded := url.QueryEscape(string(jsonData))

	http.SetCookie(w, &http.Cookie{
		Name:     "alerts",
		Value:    string(encoded),
		Path:     "/",
		Expires:  time.Now().Add(5 * time.Minute),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
}
