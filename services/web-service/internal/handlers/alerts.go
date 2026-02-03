package handlers

import (
	"net/http"

	"web-service/internal/models"
)

func SetAlerts(r *http.Request, alerts ...models.Alert) {
	a, ok := r.Context().Value("alerts").([]models.Alert)
	if !ok {
		return
	}

	a = append(a, alerts...)
}

func GetAlerts(r *http.Request) []models.Alert {
	alerts := []models.Alert{}
	a, ok := r.Context().Value("alerts").([]models.Alert)
	if !ok {
		return alerts
	}

	alerts = append(alerts, a...)

	a = nil
	
	return alerts
}
