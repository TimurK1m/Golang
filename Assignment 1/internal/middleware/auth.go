package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const apiKey = "secret12345"

// ===== API KEY MIDDLEWARE =====

func APIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-KEY")

		if key != apiKey {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "unauthorized",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timestamp := time.Now().Format(time.RFC3339)
		log.Printf(
			"%s %s %s {request received}",
			timestamp,
			r.Method,
			r.URL.Path,
		)

		next.ServeHTTP(w, r)
	})
}
