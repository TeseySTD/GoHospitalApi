package middleware

import (
	"log"
	"net/http"
	"os"
	"time"
)

const (
	LogFile       = "api_logs.txt"
	AuthHeaderKey = "X-API-Key"
	ValidAPIKey   = "my-secret-api-key-12345"
)

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		next(w, r)
		
		logEntry := time.Now().Format("2006-01-02 15:04:05") + 
			" | Method: " + r.Method + 
			" | Path: " + r.URL.Path +
			" | Query: " + r.URL.RawQuery +
			" | Duration: " + time.Since(start).String() + "\n"
		
		// Записуємо у файл
		f, err := os.OpenFile(LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Error opening log file: %v", err)
			return
		}
		defer f.Close()
		
		if _, err := f.WriteString(logEntry); err != nil {
			log.Printf("Error writing to log file: %v", err)
		}
	}
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get(AuthHeaderKey)
		
		if apiKey == "" {
			http.Error(w, `{"error": "API key is missing"}`, http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			return
		}
		
		if apiKey != ValidAPIKey {
			http.Error(w, `{"error": "Invalid API key"}`, http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			return
		}
		
		next(w, r)
	}
}

func Chain(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}