package handlers

import (
	"encoding/json"
	"net/http"

	"golang.org/x/time/rate"
)

func RateLimit(limiter *rate.Limiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			message := map[string]string{
				"message": "Rate limit exceeded, try again later!",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			err := json.NewEncoder(w).Encode(message)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
