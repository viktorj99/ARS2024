package repository

import (
	"errors"
	"fmt"
	"projekat/model"
)

type ConfigGroupInMemRepository struct {
	configGroups map[string]model.ConfigGroup
}

// Konstruktor
func NewConfigGroupInMemRepository() model.ConfigGroupRepository {
	return ConfigGroupInMemRepository{
		configGroups: make(map[string]model.ConfigGroup),
	}
}

func (c ConfigGroupInMemRepository) AddConfigGroup(configGroup model.ConfigGroup) {
	key := fmt.Sprintf("%s/%d", configGroup.Name, configGroup.Version)
	c.configGroups[key] = configGroup
}

func (c ConfigGroupInMemRepository) GetConfigGroup(name string, version int) (model.ConfigGroup, error) {
	key := fmt.Sprintf("%s/%d", name, version)
	configGroup, ok := c.configGroups[key]
	if !ok {
		return model.ConfigGroup{}, errors.New("config group not found")
	}
	return configGroup, nil
}

func (c ConfigGroupInMemRepository) DeleteConfigGroup(name string, version int) error {
	key := fmt.Sprintf("%s/%d", name, version)
	if _, exists := c.configGroups[key]; !exists {
		return fmt.Errorf("config not found")
	}
	delete(c.configGroups, key)
	return nil

}

func (c ConfigGroupInMemRepository) AddConfigToGroup(name string, version int, config model.Config) error {
	key := fmt.Sprintf("%s/%d", name, version)
	group, ok := c.configGroups[key]
	if !ok {
		return fmt.Errorf("config group not found")
	}
	group.Configurations = append(group.Configurations, config)
	c.configGroups[key] = group
	return nil
}
