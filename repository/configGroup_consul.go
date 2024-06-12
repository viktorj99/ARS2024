package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"projekat/model"
	"strings"

	"github.com/hashicorp/consul/api"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type ConfigGroupConsulRepository struct {
	client *api.Client
	tracer trace.Tracer
}

// Constructor to create a new ConfigGroupConsulRepository
func NewConfigGroupConsulRepository() (*ConfigGroupConsulRepository, error) {
	consulAddress := fmt.Sprintf("%s:%s", os.Getenv("DB"), os.Getenv("DBPORT"))
	config := api.DefaultConfig()
	config.Address = consulAddress

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	tracer := otel.Tracer("ConfigGroupConsulRepository")

	return &ConfigGroupConsulRepository{
		client: client,
		tracer: tracer,
	}, nil
}

func (r *ConfigGroupConsulRepository) AddConfigGroup(ctx context.Context, configGroup model.ConfigGroup) error {
	ctx, span := r.tracer.Start(ctx, "AddConfigGroup")
	defer span.End()

	kv := r.client.KV()
	data, err := json.Marshal(configGroup)
	if err != nil {
		return err
	}
	p := &api.KVPair{
		Key:   fmt.Sprintf("configgroup/%s/%d", configGroup.Name, configGroup.Version),
		Value: data,
	}
	_, err = kv.Put(p, nil)
	return err
}

func (r *ConfigGroupConsulRepository) GetConfigGroup(ctx context.Context, name string, version int) (model.ConfigGroup, error) {
	ctx, span := r.tracer.Start(ctx, "GetConfigGroup")
	defer span.End()

	kv := r.client.KV()
	key := fmt.Sprintf("configgroup/%s/%d", name, version)
	pairs, _, err := kv.List(key, nil)
	if err != nil {
		return model.ConfigGroup{}, err
	}
	if len(pairs) == 0 {
		return model.ConfigGroup{}, fmt.Errorf("config group not found")
	}
	var configGroup model.ConfigGroup
	err = json.Unmarshal(pairs[0].Value, &configGroup)
	if err != nil {
		return model.ConfigGroup{}, err
	}
	return configGroup, nil
}

func (r *ConfigGroupConsulRepository) DeleteConfigGroup(ctx context.Context, name string, version int) error {
	ctx, span := r.tracer.Start(ctx, "DeleteConfigGroup")
	defer span.End()

	kv := r.client.KV()
	key := fmt.Sprintf("configgroup/%s/%d", name, version)
	pairs, _, err := kv.List(key, nil)
	if err != nil {
		return fmt.Errorf("error checking config group: %v", err)
	}
	if len(pairs) == 0 {
		return fmt.Errorf("config group not found")
	}
	for _, pair := range pairs {
		_, err = kv.Delete(pair.Key, nil)
		if err != nil {
			return fmt.Errorf("error deleting config group: %v", err)
		}
	}
	return nil
}

func (r *ConfigGroupConsulRepository) AddConfigToGroup(ctx context.Context, name string, version int, config model.Config) error {
	ctx, span := r.tracer.Start(ctx, "AddConfigToGroup")
	defer span.End()

	configGroup, err := r.GetConfigGroup(ctx, name, version)
	if err != nil {
		return err
	}
	configGroup.Configurations = append(configGroup.Configurations, config)
	return r.AddConfigGroup(ctx, configGroup)
}

func (r *ConfigGroupConsulRepository) DeleteConfigFromGroup(ctx context.Context, groupName string, groupVersion int, configName string, configVersion int) error {
	ctx, span := r.tracer.Start(ctx, "DeleteConfigFromGroup")
	defer span.End()

	configGroup, err := r.GetConfigGroup(ctx, groupName, groupVersion)
	if err != nil {
		return err
	}
	for i, config := range configGroup.Configurations {
		if config.Name == configName && config.Version == configVersion {
			configGroup.Configurations = append(configGroup.Configurations[:i], configGroup.Configurations[i+1:]...)
			break
		}
	}
	return r.AddConfigGroup(ctx, configGroup)
}

func (r *ConfigGroupConsulRepository) GetConfigsFromGroupByLabel(ctx context.Context, groupName string, groupVersion int, labels string) ([]model.Config, error) {
	ctx, span := r.tracer.Start(ctx, "GetConfigsFromGroupByLabel")
	defer span.End()

	labelMap := make(map[string]string)
	labelPairs := strings.Split(labels, ";")

	for _, pair := range labelPairs {
		parts := strings.Split(pair, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid label format")
		}
		labelMap[parts[0]] = parts[1]
	}

	configGroup, err := r.GetConfigGroup(ctx, groupName, groupVersion)
	if err != nil {
		return nil, err
	}

	var configs []model.Config
	for _, config := range configGroup.Configurations {
		if labelMapsAreEqual(labelMap, config.Labels) {
			configs = append(configs, config)
		}
	}

	if len(configs) == 0 {
		return nil, fmt.Errorf("config not found")
	}

	return configs, nil
}

func (r *ConfigGroupConsulRepository) DeleteConfigsFromGroupByLabel(ctx context.Context, groupName string, groupVersion int, labels string) error {
	ctx, span := r.tracer.Start(ctx, "DeleteConfigsFromGroupByLabel")
	defer span.End()

	labelMap := make(map[string]string)
	labelPairs := strings.Split(labels, ";")

	for _, pair := range labelPairs {
		parts := strings.Split(pair, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid label format")
		}
		labelMap[parts[0]] = parts[1]
	}

	configGroup, err := r.GetConfigGroup(ctx, groupName, groupVersion)
	if err != nil {
		return err
	}

	var updatedConfigs []model.Config
	found := false
	for _, config := range configGroup.Configurations {
		if labelMapsAreEqual(labelMap, config.Labels) {
			found = true
		} else {
			updatedConfigs = append(updatedConfigs, config)
		}
	}

	if !found {
		return fmt.Errorf("config not found")
	}

	configGroup.Configurations = updatedConfigs
	return r.AddConfigGroup(ctx, configGroup)
}

func labelMapsAreEqual(map1, map2 map[string]string) bool {
	if len(map1) != len(map2) {
		return false
	}
	for key, value := range map1 {
		if map2[key] != value {
			return false
		}
	}
	return true
}
