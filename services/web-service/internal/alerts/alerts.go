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

	AlertKey    string = "sc_alerts"
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

		}
	}

	log.Printf("Read alerts: %v", alerts)
	return context.WithValue(r.Context(), AlertKey, alerts)
}

func WriteAlerts(w http.ResponseWriter, r *http.Request) {
	alerts, ok := r.Context().Value(AlertKey).([]Alert)
	if !ok {
		return
	}

	log.Printf("Write alerts: %v", alerts)

	if len(alerts) == 0 {
		http.SetCookie(w, &http.Cookie{
			Name:     AlertKey,
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})
		return
	}

	jsonData, _ := json.Marshal(alerts)

	encoded := url.QueryEscape(string(jsonData))

	http.SetCookie(w, &http.Cookie{
		Name:     AlertKey,
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
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

func SetAlert(r *http.Request, alert Alert) {
	alerts := GetAlerts(r.Context())
	alerts = append(alerts, alert)

	log.Printf("Set alerts: %v", alerts)

	ctx := context.WithValue(r.Context(), AlertKey, alerts)
	*r = *r.WithContext(ctx)
}

func FlashSuccess(w http.ResponseWriter, r *http.Request, msg string) {
	SetAlert(r, Alert{
		Type: AlertSuccess,
		Msg:  msg,
	})
}

func FlashInfo(w http.ResponseWriter, r *http.Request, msg string) {
	SetAlert(r, Alert{
		Type: AlertInfo,
		Msg:  msg,
	})
}

func FlashWarning(w http.ResponseWriter, r *http.Request, msg string) {
	SetAlert(r, Alert{
		Type: AlertWarning,
		Msg:  msg,
	})
}

func FlashError(w http.ResponseWriter, r *http.Request, msg string) {
	SetAlert(r, Alert{
		Type: AlertError,
		Msg:  msg,
	})
}
