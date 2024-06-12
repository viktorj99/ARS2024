package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/consul/api"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type IdempotencyRepository struct {
	client *api.Client
	tracer trace.Tracer
}

type IdempotencyData struct {
	BodyHash  string    `json:"body_hash"`
	Timestamp time.Time `json:"timestamp"`
}

func NewIdempotencyRepository() (*IdempotencyRepository, error) {
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", os.Getenv("DB"), os.Getenv("DBPORT"))

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	tracer := otel.Tracer("IdempotencyRepository")

	return &IdempotencyRepository{
		client: client,
		tracer: tracer,
	}, nil
}

func (r *IdempotencyRepository) GetIdempotencyKey(ctx context.Context, endpoint, idempotencyKey string) (string, error) {
	ctx, span := r.tracer.Start(ctx, "GetIdempotencyKey")
	defer span.End()

	kv := r.client.KV()
	key := fmt.Sprintf("idempotency%s/%s", endpoint, idempotencyKey)
	log.Printf("Getting idempotency key: %s", key)
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		log.Printf("Error getting idempotency key: %v", err)
		return "", err
	}
	if pair == nil {
		log.Printf("Idempotency key not found: %s", key)
		return "", nil
	}

	var data IdempotencyData
	err = json.Unmarshal(pair.Value, &data)
	if err != nil {
		log.Printf("Error unmarshaling idempotency data: %v", err)
		return "", err
	}
	log.Printf("Idempotency data found: %v", data)
	return data.BodyHash, nil
}

func (r *IdempotencyRepository) SaveIdempotencyKey(ctx context.Context, endpoint, idempotencyKey, bodyHash string, timestamp time.Time) error {
	ctx, span := r.tracer.Start(ctx, "SaveIdempotencyKey")
	defer span.End()

	kv := r.client.KV()
	data := IdempotencyData{
		BodyHash:  bodyHash,
		Timestamp: timestamp,
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling idempotency data: %v", err)
		return err
	}
	p := &api.KVPair{
		Key:   fmt.Sprintf("idempotency%s/%s", endpoint, idempotencyKey),
		Value: dataBytes,
	}
	log.Printf("Saving idempotency key: %s with data: %v", p.Key, data)
	_, err = kv.Put(p, nil)
	if err != nil {
		log.Printf("Error saving idempotency key: %v", err)
	}
	return err
}
