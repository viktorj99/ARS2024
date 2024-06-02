package handlers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"sync"
	"time"
)

var (
	idempotencyData   = make(map[string]map[string]time.Time)
	idempotencyBodies = make(map[string]map[string]string)
	mutex             = &sync.Mutex{}
)

func IdempotencyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			next.ServeHTTP(w, r)
			return
		}

		idempotencyKey := r.Header.Get("Idempotency-Key")
		if idempotencyKey == "" {
			http.Error(w, "Idempotency-Key header missing", http.StatusBadRequest)
			return
		}

		endpoint := r.URL.Path

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		bodyHash := hashBody(bodyBytes)

		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		mutex.Lock()
		defer mutex.Unlock()

		if _, exists := idempotencyData[endpoint]; !exists {
			idempotencyData[endpoint] = make(map[string]time.Time)
			idempotencyBodies[endpoint] = make(map[string]string)
		}

		if existingHash, exists := idempotencyBodies[endpoint][idempotencyKey]; exists {
			if existingHash == bodyHash {
				http.Error(w, "Duplicate request with same body", http.StatusConflict)
				return
			}
		}

		idempotencyData[endpoint][idempotencyKey] = time.Now()
		idempotencyBodies[endpoint][idempotencyKey] = bodyHash

		next.ServeHTTP(w, r)
	})
}

func hashBody(body []byte) string {
	hash := sha256.Sum256(body)
	return hex.EncodeToString(hash[:])
}
