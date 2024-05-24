package repository

import (
	"encoding/json"
	"fmt"
	"projekat/model"

	"github.com/hashicorp/consul/api"
)

type ConfigGroupConsulRepository struct {
	client *api.Client
}

// Constructor to create a new ConfigGroupConsulRepository
func NewConfigGroupConsulRepository() (*ConfigGroupConsulRepository, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}
	return &ConfigGroupConsulRepository{client: client}, nil
}

func (r *ConfigGroupConsulRepository) AddConfigGroup(configGroup model.ConfigGroup) error {
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

func (r *ConfigGroupConsulRepository) GetConfigGroup(name string, version int) (model.ConfigGroup, error) {
	kv := r.client.KV()
	key := fmt.Sprintf("configgroup/%s/%d", name, version)
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return model.ConfigGroup{}, err
	}
	if pair == nil {
		return model.ConfigGroup{}, fmt.Errorf("config group not found")
	}
	var configGroup model.ConfigGroup
	err = json.Unmarshal(pair.Value, &configGroup)
	if err != nil {
		return model.ConfigGroup{}, err
	}
	return configGroup, nil
}

func (r *ConfigGroupConsulRepository) DeleteConfigGroup(name string, version int) error {
	kv := r.client.KV()
	key := fmt.Sprintf("configgroup/%s/%d", name, version)
	_, err := kv.Delete(key, nil)
	return err
}

func (r *ConfigGroupConsulRepository) AddConfigToGroup(name string, version int, config model.Config) error {
	configGroup, err := r.GetConfigGroup(name, version)
	if err != nil {
		return err
	}
	configGroup.Configurations = append(configGroup.Configurations, config)
	return r.AddConfigGroup(configGroup)
}

func (r *ConfigGroupConsulRepository) DeleteConfigFromGroup(groupName string, groupVersion int, configName string, configVersion int) error {
	configGroup, err := r.GetConfigGroup(groupName, groupVersion)
	if err != nil {
		return err
	}
	for i, config := range configGroup.Configurations {
		if config.Name == configName && config.Version == configVersion {
			configGroup.Configurations = append(configGroup.Configurations[:i], configGroup.Configurations[i+1:]...)
			break
		}
	}
	return r.AddConfigGroup(configGroup)
}

func (r *ConfigGroupConsulRepository) GetConfigsFromGroupByLabel(groupName string, groupVersion int, labels string) ([]model.Config, error) {
	configGroup, err := r.GetConfigGroup(groupName, groupVersion)
	if err != nil {
		return nil, err
	}
	var configs []model.Config
	for _, config := range configGroup.Configurations {
		for k := range config.Labels {
			if k == labels {
				configs = append(configs, config)
				break
			}
		}
	}
	return configs, nil
}

func (r *ConfigGroupConsulRepository) DeleteConfigsFromGroupByLabel(groupName string, groupVersion int, labels string) error {
	configGroup, err := r.GetConfigGroup(groupName, groupVersion)
	if err != nil {
		return err
	}
	var configs []model.Config
	for _, config := range configGroup.Configurations {
		keep := true
		for k := range config.Labels {
			if k == labels {
				keep = false
				break
			}
		}
		if keep {
			configs = append(configs, config)
		}
	}
	configGroup.Configurations = configs
	return r.AddConfigGroup(configGroup)
}
