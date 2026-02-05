package alerts

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

type AlertType string

const (
	AlertSuccess AlertType = "success"
	AlertError   AlertType = "danger"
	AlertWarning AlertType = "warning"
	AlertInfo    AlertType = "info"

	AlertKey    string = "alerts"
	AlertExpire int    = 300
)

type Alert struct {
	Type AlertType
	Msg  string
}

func ReadAlerts(w http.ResponseWriter, r *http.Request) context.Context {
	alerts := []Alert{}
	if r != nil {
		cookie, err := r.Cookie(AlertKey)
		if err == nil {
			decoded, _ := url.QueryUnescape(cookie.Value)

			json.Unmarshal([]byte(decoded), &alerts)

			log.Printf("Read alerts: %v", alerts)

			http.SetCookie(w, &http.Cookie{
				Name:     AlertKey,
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
			})
		}
	}

	return context.WithValue(r.Context(), AlertKey, alerts)
}

func GetAlerts(ctx context.Context) []Alert {
	alerts, ok := ctx.Value(AlertKey).([]Alert)
	if !ok {
		alerts = []Alert{}
		log.Printf("Failed to get alerts!")
	}

	log.Printf("Get alerts: %v", alerts)

	return alerts
}

func SetAlert(ctx context.Context, w http.ResponseWriter, alert Alert) {
	alerts := GetAlerts(ctx)
	alerts = append(alerts, alert)

	log.Printf("Set alerts: %v", alerts)

	jsonData, _ := json.Marshal(alerts)

	encoded := url.QueryEscape(string(jsonData))

	http.SetCookie(w, &http.Cookie{
		Name:     AlertKey,
		Value:    string(encoded),
		Path:     "/",
		MaxAge:   AlertExpire,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
}

func FlashSuccess(w http.ResponseWriter, r *http.Request, msg string) {
	SetAlert(r.Context(), w, Alert{
		Type: AlertSuccess,
		Msg:  msg,
	})
}

func FlashInfo(w http.ResponseWriter, r *http.Request, msg string) {
	SetAlert(r.Context(), w, Alert{
		Type: AlertInfo,
		Msg:  msg,
	})
}

func FlashWarning(w http.ResponseWriter, r *http.Request, msg string) {
	SetAlert(r.Context(), w, Alert{
		Type: AlertWarning,
		Msg:  msg,
	})
}

func FlashError(w http.ResponseWriter, r *http.Request, msg string) {
	SetAlert(r.Context(), w, Alert{
		Type: AlertError,
		Msg:  msg,
	})
}
