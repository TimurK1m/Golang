package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetRate(t *testing.T) {
	
	t.Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"base": "USD", "target": "EUR", "rate": 0.92}`))
		}))
		defer server.Close() 

		service := NewExchangeService(server.URL)
		rate, err := service.GetRate("USD", "EUR")

		assert.NoError(t, err)
		assert.Equal(t, 0.92, rate)
	})

	
	t.Run("API Business Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "invalid currency pair"}`))
		}))
		defer server.Close()

		service := NewExchangeService(server.URL)
		_, err := service.GetRate("INVALID", "EUR")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "api error: invalid currency pair")
	})

	
	t.Run("Malformed JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{ "invalid_json": `)) 
		}))
		defer server.Close()

		service := NewExchangeService(server.URL)
		_, err := service.GetRate("USD", "EUR")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "decode error")
	})

	
	t.Run("Timeout Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(6 * time.Second) 
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		service := NewExchangeService(server.URL)
		_, err := service.GetRate("USD", "EUR")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "network error")
	})

	
	t.Run("Server Error 500", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			
			w.Write([]byte(`{}`)) 
		}))
		defer server.Close()

		service := NewExchangeService(server.URL)
		_, err := service.GetRate("USD", "EUR")

		assert.Error(t, err)
		
		assert.Contains(t, err.Error(), "unexpected status: 500")
	})

	
	t.Run("Empty Body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			
		}))
		defer server.Close()

		service := NewExchangeService(server.URL)
		_, err := service.GetRate("USD", "EUR")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "decode error")
	})
}