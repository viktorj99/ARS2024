package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"projekat/model"

	"github.com/hashicorp/consul/api"
)

type ConfigConsulRepository struct {
	client *api.Client
}

func NewConfigConsulRepository() (*ConfigConsulRepository, error) {
	consulAddress := fmt.Sprintf("%s:%s", os.Getenv("DB"), os.Getenv("DBPORT"))
	config := api.DefaultConfig()
	config.Address = consulAddress

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &ConfigConsulRepository{client: client}, nil
}

func (r *ConfigConsulRepository) AddConfig(config model.Config) error {
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

func (r *ConfigConsulRepository) GetConfig(name string, version int) (model.Config, error) {
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

func (r *ConfigConsulRepository) DeleteConfig(name string, version int) error {
	kv := r.client.KV()
	key := fmt.Sprintf("config/%s/%d", name, version)
	_, err := kv.Delete(key, nil)
	return err
}
