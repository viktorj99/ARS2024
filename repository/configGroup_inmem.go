package repository

// import (
// 	"errors"
// 	"fmt"
// 	"projekat/model"
// 	"strings"
// )

// type ConfigGroupInMemRepository struct {
// 	configGroups map[string]model.ConfigGroup
// }

// // Konstruktor
// func NewConfigGroupInMemRepository() model.ConfigGroupRepository {
// 	return ConfigGroupInMemRepository{
// 		configGroups: make(map[string]model.ConfigGroup),
// 	}
// }

// func (c ConfigGroupInMemRepository) AddConfigGroup(configGroup model.ConfigGroup) {
// 	key := fmt.Sprintf("%s/%d", configGroup.Name, configGroup.Version)
// 	c.configGroups[key] = configGroup
// }

// func (c ConfigGroupInMemRepository) GetConfigGroup(name string, version int) (model.ConfigGroup, error) {
// 	key := fmt.Sprintf("%s/%d", name, version)
// 	configGroup, ok := c.configGroups[key]
// 	if !ok {
// 		return model.ConfigGroup{}, errors.New("config group not found")
// 	}
// 	return configGroup, nil
// }

// func (c ConfigGroupInMemRepository) DeleteConfigGroup(name string, version int) error {
// 	key := fmt.Sprintf("%s/%d", name, version)
// 	if _, exists := c.configGroups[key]; !exists {
// 		return fmt.Errorf("config not found")
// 	}
// 	delete(c.configGroups, key)
// 	return nil

// }

// func (c ConfigGroupInMemRepository) AddConfigToGroup(name string, version int, config model.Config) error {
// 	key := fmt.Sprintf("%s/%d", name, version)
// 	group, ok := c.configGroups[key]
// 	if !ok {
// 		return fmt.Errorf("config group not found")
// 	}
// 	group.Configurations = append(group.Configurations, config)
// 	c.configGroups[key] = group
// 	return nil
// }

// func (c ConfigGroupInMemRepository) DeleteConfigFromGroup(groupName string, groupVersion int, configName string, configVersion int) error {
// 	groupKey := fmt.Sprintf("%s/%d", groupName, groupVersion)
// 	group, ok := c.configGroups[groupKey]
// 	if !ok {
// 		return fmt.Errorf("config group not found")
// 	}

// 	// Check if config exists in group.Configurations
// 	configFound := false
// 	for i, config := range group.Configurations {
// 		if config.Name == configName && config.Version == configVersion {
// 			// Remove the config from group.Configurations
// 			group.Configurations = append(group.Configurations[:i], group.Configurations[i+1:]...)
// 			configFound = true
// 			break
// 		}
// 	}

// 	if !configFound {
// 		return fmt.Errorf("config not found")
// 	}

// 	c.configGroups[groupKey] = group
// 	return nil
// }

// func (c ConfigGroupInMemRepository) GetConfigsFromGroupByLabel(groupName string, groupVersion int, labels string) ([]model.Config, error) {
// 	labelMap := make(map[string]string)
// 	labelPairs := strings.Split(labels, ";")

// 	for _, pair := range labelPairs {
// 		parts := strings.Split(pair, ":")
// 		if len(parts) != 2 {
// 			return nil, errors.New("invalid label format")
// 		}
// 		labelMap[parts[0]] = parts[1]
// 	}

// 	key := fmt.Sprintf("%s/%d", groupName, groupVersion)
// 	configGroup, ok := c.configGroups[key]
// 	if !ok {
// 		return nil, errors.New("config group not found")
// 	}

// 	var result []model.Config

// 	for _, config := range configGroup.Configurations {
// 		if labelMapsAreEqual(labelMap, config.Labels) {
// 			result = append(result, config)
// 		}
// 	}

// 	if result == nil {
// 		return nil, errors.New("config not found")
// 	}

// 	return result, nil
// }

// func (c ConfigGroupInMemRepository) DeleteConfigsFromGroupByLabel(groupName string, groupVersion int, labels string) error {
// 	labelMap := make(map[string]string)
// 	labelPairs := strings.Split(labels, ";")
// 	for _, pair := range labelPairs {
// 		parts := strings.Split(pair, ":")
// 		if len(parts) != 2 {
// 			return errors.New("invalid label format")
// 		}
// 		labelMap[parts[0]] = parts[1]
// 	}

// 	groupKey := fmt.Sprintf("%s/%d", groupName, groupVersion)
// 	group, ok := c.configGroups[groupKey]
// 	if !ok {
// 		return errors.New("config group not found")
// 	}

// 	found := false
// 	for _, config := range group.Configurations {
// 		if labelMapsAreEqual(labelMap, config.Labels) {
// 			err := c.DeleteConfigFromGroup(groupName, groupVersion, config.Name, config.Version)
// 			if err != nil {
// 				return err
// 			}
// 			found = true
// 		}
// 	}

// 	if !found {
// 		return errors.New("config not found")
// 	}

// 	return nil
// }

// func labelMapsAreEqual(map1, map2 map[string]string) bool {
// 	if len(map1) != len(map2) {
// 		return false
// 	}
// 	for key, value := range map1 {
// 		if map2[key] != value {
// 			return false
// 		}
// 	}
// 	return true
// }
