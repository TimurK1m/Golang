package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"task-1/internal"
	"time"
)

func main() {
	attempt := 0

	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt++

		if attempt <= 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Println("Server: returning 503")
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
		fmt.Println("Server: returning 200")
	}))
	defer server.Close()

	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := &http.Client{}

	err := internal.ExecutePayment(ctx, client, server.URL)
	if err != nil {
		fmt.Println("Final error:", err)
	}
}