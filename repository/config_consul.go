package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"projekat/model"

	"github.com/hashicorp/consul/api"
	"go.opentelemetry.io/otel/trace"
)

type ConfigConsulRepository struct {
	client *api.Client
	tracer trace.Tracer
}

func NewConfigConsulRepository(tracer trace.Tracer) (*ConfigConsulRepository, error) {
	consulAddress := fmt.Sprintf("%s:%s", os.Getenv("DB"), os.Getenv("DBPORT"))
	config := api.DefaultConfig()
	config.Address = consulAddress

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &ConfigConsulRepository{client: client, tracer: tracer}, nil
}

func (r *ConfigConsulRepository) AddConfig(ctx context.Context, config model.Config) error {
	ctx, span := r.tracer.Start(ctx, "AddConfig")
	defer span.End()

	kv := r.client.KV()
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	p := &api.KVPair{
		Key:   fmt.Sprintf("config/%s/%d", config.Name, config.Version),
		Value: data,
	}
	_, err = kv.Put(p, nil)
	return err
}

func (r *ConfigConsulRepository) GetConfig(ctx context.Context, name string, version int) (model.Config, error) {
	ctx, span := r.tracer.Start(ctx, "GetConfig")
	defer span.End()

	kv := r.client.KV()
	key := fmt.Sprintf("config/%s/%d", name, version)
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return model.Config{}, err
	}
	if pair == nil {
		return model.Config{}, fmt.Errorf("config not found")
	}
	var config model.Config
	err = json.Unmarshal(pair.Value, &config)
	if err != nil {
		return model.Config{}, err
	}
	return config, nil
}

func (r *ConfigConsulRepository) DeleteConfig(ctx context.Context, name string, version int) error {
	ctx, span := r.tracer.Start(ctx, "DeleteConfig")
	defer span.End()

	kv := r.client.KV()
	key := fmt.Sprintf("config/%s/%d", name, version)

	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return fmt.Errorf("error checking config: %v", err)
	}
	if pair == nil {
		return fmt.Errorf("config not found")
	}

	_, err = kv.Delete(key, nil)
	if err != nil {
		return fmt.Errorf("error deleting config : %v", err)
	}
	return nil
}
