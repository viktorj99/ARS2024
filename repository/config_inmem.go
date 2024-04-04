package repository

import "projekat/model"

type ConfigInMemRepository struct {
	configs map[string]model.Config
}

func (c ConfigInMemRepository) AddConfig(config model.Config) error {
	//TODO implement me
	panic("implement me")
}

func (c ConfigInMemRepository) DeleteConfig(name string, version string) error {
	//TODO implement me
	panic("implement me")
}

func (c ConfigInMemRepository) UpdateConfig(name string, version string, config model.Config) error {
	//TODO implement me
	panic("implement me")
}

func (c ConfigInMemRepository) ListConfigs() ([]model.Config, error) {
	//TODO implement me
	panic("implement me")
}

func (c ConfigInMemRepository) FindConfigByNameAndVersion(name string, version string) (*model.Config, error) {
	//TODO implement me
	panic("implement me")
}

func NewConfigInMemRepository() model.ConfigRepository {
	return ConfigInMemRepository{
		configs: make(map[string]model.Config),
	}
}

func (c ConfigInMemRepository) Get() {

}
