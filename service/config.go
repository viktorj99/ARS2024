package service

import (
	"projekat/model"
)

//Klasa
type ConfigService struct {
	repository model.ConfigRepository
}

func NewConfigService(repository model.ConfigRepository) ConfigService {
	return ConfigService {
		repository: repository,
	}
}

func (s ConfigService) AddConfig(config model.Config) {
	s.repository.AddConfig(config);
}

func (s ConfigService) GetConfig (name string, version int) (model.Config, error) {
	return s.repository.GetConfig(name, version);
}

func (s ConfigService) DeleteConfig (name string, version int) error {
	return s.repository.DeleteConfig(name, version);
}