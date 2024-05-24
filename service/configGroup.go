package service

import (
	"projekat/model"
)

// Klasa
type ConfigGroupService struct {
	repository model.ConfigGroupRepository
}

func NewConfigGroupService(repository model.ConfigGroupRepository) ConfigGroupService {
	return ConfigGroupService{
		repository: repository,
	}
}

func (s ConfigGroupService) AddConfigGroup(config model.ConfigGroup) error {
	return s.repository.AddConfigGroup(config)
}

func (s ConfigGroupService) GetConfigGroup(name string, version int) (model.ConfigGroup, error) {
	return s.repository.GetConfigGroup(name, version)
}

func (s ConfigGroupService) DeleteConfigGroup(name string, version int) error {
	return s.repository.DeleteConfigGroup(name, version)
}

func (s ConfigGroupService) AddConfigToGroup(name string, version int, config model.Config) error {
	return s.repository.AddConfigToGroup(name, version, config)
}

func (s ConfigGroupService) DeleteConfigFromGroup(groupName string, groupVersion int, configName string, configVersion int) error {
	return s.repository.DeleteConfigFromGroup(groupName, groupVersion, configName, configVersion)
}

func (s ConfigGroupService) GetConfigsFromGroupByLabel(groupName string, groupVersion int, labels string) ([]model.Config, error) {
	return s.repository.GetConfigsFromGroupByLabel(groupName, groupVersion, labels)
}

func (s ConfigGroupService) DeleteConfigsFromGroupByLabel(groupName string, groupVersion int, labels string) error {
	return s.repository.DeleteConfigsFromGroupByLabel(groupName, groupVersion, labels)
}
