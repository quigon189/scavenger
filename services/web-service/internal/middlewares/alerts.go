package middlewares

import (
	"log"
	"net/http"
	"web-service/internal/alerts"
)

func ProccessAlerts(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := alerts.ReadAlerts(w, r)

		*r = *r.WithContext(ctx)

		next.ServeHTTP(w, r)

		alerts.WriteAlerts(w,r)

		log.Printf("Headers: %v", w.Header())
	})
}
