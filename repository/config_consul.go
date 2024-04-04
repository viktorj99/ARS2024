package repository

import "projekat/model"

type ConfigConsulRepository struct {
}

func (c ConfigConsulRepository) AddConfig(config model.Config) error {
	//TODO implement me
	panic("implement me")
}

func (c ConfigConsulRepository) DeleteConfig(name string, version string) error {
	//TODO implement me
	panic("implement me")
}

func (c ConfigConsulRepository) UpdateConfig(name string, version string, config model.Config) error {
	//TODO implement me
	panic("implement me")
}

func (c ConfigConsulRepository) ListConfigs() ([]model.Config, error) {
	//TODO implement me
	panic("implement me")
}

func (c ConfigConsulRepository) FindConfigByNameAndVersion(name string, version string) (*model.Config, error) {
	//TODO implement me
	panic("implement me")
}

func NewConfigConsulRepository() model.ConfigRepository {
	return ConfigConsulRepository{}
}
