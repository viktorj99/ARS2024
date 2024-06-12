package handlers

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"projekat/repository"
	"sync"
	"time"
)

var (
	idempotencyRepository *repository.IdempotencyRepository
	mutex                 = &sync.Mutex{}
)

func SetIdempotencyRepository(repo *repository.IdempotencyRepository) {
	idempotencyRepository = repo
}

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
		fmt.Println(endpoint)

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		bodyHash := hashBody(bodyBytes)

		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		mutex.Lock()
		defer mutex.Unlock()

		existingHash, err := idempotencyRepository.GetIdempotencyKey(context.Background(), endpoint, idempotencyKey)
		if err != nil {
			http.Error(w, "Error checking idempotency key", http.StatusInternalServerError)
			return
		}

		if existingHash == bodyHash {
			log.Printf("Duplicate request detected: %s for endpoint: %s", idempotencyKey, endpoint)
			http.Error(w, "Duplicate request with same body.", http.StatusConflict)
			return
		}

		err = idempotencyRepository.SaveIdempotencyKey(context.Background(), endpoint, idempotencyKey, bodyHash, time.Now())
		if err != nil {
			http.Error(w, "Error saving idempotency key", http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func hashBody(body []byte) string {
	hash := sha256.Sum256(body)
	return hex.EncodeToString(hash[:])
}
