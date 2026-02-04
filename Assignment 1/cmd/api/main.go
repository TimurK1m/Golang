package main

import (
	"log"
	"net/http"

	"Assignment-1/internal/handlers"
	"Assignment-1/internal/middleware"
)

func main() {
	mux := http.NewServeMux()

	taskHandler := http.HandlerFunc(handlers.TasksHandler)


	// middleware chain
	handler := middleware.Logger(
	middleware.APIKey(taskHandler),
)


	mux.Handle("/tasks", handler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
