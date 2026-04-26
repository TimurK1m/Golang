package internal

import (
	"bytes"
	"fmt"
	"net/http"
)

func IdempotencyMiddleware(store *Store, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		key := r.Header.Get("Idempotency-Key")
		if key == "" {
			http.Error(w, "Missing Idempotency-Key", http.StatusBadRequest)
			return
		}

		rec, exists := store.Get(key)

		if exists {
			switch rec.Status {
			case StatusProcessing:
				http.Error(w, "Request already processing", http.StatusConflict)
				return

			case StatusCompleted:
				fmt.Println("Returning cached response")
				w.WriteHeader(rec.StatusCode)
				w.Write(rec.Body)
				return
			}
		}

		
		store.SetProcessing(key)
		fmt.Println("Processing started")

		
		recorder := &responseRecorder{
			ResponseWriter: w,
			body:           &bytes.Buffer{},
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		
		store.SetCompleted(key, recorder.statusCode, recorder.body.Bytes())
	})
}

type responseRecorder struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}