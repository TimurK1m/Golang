package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"task-2/internal"
	"time"

	"github.com/google/uuid"
)

func paymentHandler(w http.ResponseWriter, r *http.Request) {
	
	time.Sleep(2 * time.Second)

	resp := map[string]interface{}{
		"status":         "paid",
		"amount":         1000,
		"transaction_id": uuid.New().String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	store := internal.NewStore()

	mux := http.NewServeMux()
	mux.Handle("/pay", internal.IdempotencyMiddleware(store, http.HandlerFunc(paymentHandler)))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		fmt.Println("Server started on :8080")
		server.ListenAndServe()
	}()

	time.Sleep(1 * time.Second)

	
	var wg sync.WaitGroup
	client := &http.Client{}
	key := "fixed-key-123"

	for i := 0; i < 7; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			req, _ := http.NewRequest("POST", "http://localhost:8080/pay", nil)
			req.Header.Set("Idempotency-Key", key)

			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Request error:", err)
				return
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)

			fmt.Printf("Request %d → Status: %d, Body: %s\n", i, resp.StatusCode, string(body))
		}(i)
	}

	wg.Wait()

	fmt.Println("\n--- Sending one more request after completion ---")

	
	req, _ := http.NewRequest("POST", "http://localhost:8080/pay", nil)
	req.Header.Set("Idempotency-Key", key)

	resp, _ := client.Do(req)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Final request → Status: %d, Body: %s\n", resp.StatusCode, string(body))
}