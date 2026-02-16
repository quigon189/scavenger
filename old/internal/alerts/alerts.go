package alerts

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

type AlertType string

type contextKey string

const AlertKey contextKey = "alerts"

const (
	AlertSuccess AlertType = "success"
	AlertError   AlertType = "danger"
	AlertWarning AlertType = "warning"
	AlertInfo    AlertType = "info"
)

type Alert struct {
	Type    AlertType `json:"type"`
	Message string    `json:"message"`
}

func SetAlert(w http.ResponseWriter, r *http.Request, alertType AlertType, message string) {
	alerts := GetAlertsFromRequest(r)
	alerts = append(alerts, Alert{Type: alertType, Message: message})

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

func GetAlerts(w http.ResponseWriter, r *http.Request) []Alert {
	alerts := GetAlertsFromRequest(r)

	http.SetCookie(w, &http.Cookie{
		Name:     "alerts",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Minute),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	return alerts
}

func GetAlertsFromRequest(r *http.Request) []Alert {
	var alerts []Alert
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

func GetAlertsFromContext(ctx context.Context) []Alert {
	if alerts, ok := ctx.Value(AlertKey).([]Alert); ok {
		return alerts
	}
	return []Alert{}
}

func FlashError(w http.ResponseWriter, r *http.Request, message string) {
	SetAlert(w, r, AlertError, message)
}

func FlashInfo(w http.ResponseWriter, r *http.Request, message string) {
	SetAlert(w, r, AlertInfo, message)
}

func FlashWarning(w http.ResponseWriter, r *http.Request, message string) {
	SetAlert(w, r, AlertWarning, message)
}

func FlashSuccess(w http.ResponseWriter, r *http.Request, message string) {
	SetAlert(w, r, AlertSuccess, message)
}
