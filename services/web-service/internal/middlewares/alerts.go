package middlewares

import (
	"net/http"
	"web-service/internal/alerts"
)

func ProccessAlerts(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := alerts.ReadAlerts(w, r)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
