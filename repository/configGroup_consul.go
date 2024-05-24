package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"projekat/model"
	"strings"

	"github.com/hashicorp/consul/api"
)

type ConfigGroupConsulRepository struct {
	client *api.Client
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
	labelMap := make(map[string]string)
	labelPairs := strings.Split(labels, ";")

	for _, pair := range labelPairs {
		parts := strings.Split(pair, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid label format")
		}
		labelMap[parts[0]] = parts[1]
	}

	configGroup, err := r.GetConfigGroup(groupName, groupVersion)
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

func (r *ConfigGroupConsulRepository) DeleteConfigsFromGroupByLabel(groupName string, groupVersion int, labels string) error {
	labelMap := make(map[string]string)
	labelPairs := strings.Split(labels, ";")

	for _, pair := range labelPairs {
		parts := strings.Split(pair, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid label format")
		}
		labelMap[parts[0]] = parts[1]
	}

	configGroup, err := r.GetConfigGroup(groupName, groupVersion)
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
	return r.AddConfigGroup(configGroup)
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
