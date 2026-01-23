package handlers

import (
	"log"
	"net/http"
	"time"
)

func LogginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		t := time.Since(start)

		log.Printf("INFO %s %s", r.URL.Path, t.String())
	})
}
