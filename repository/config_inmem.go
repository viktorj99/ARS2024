package repository

import (
	"errors"
	"fmt"
	"projekat/model"
)

//Klasa
type ConfigInMemRepository struct {
	configs map[string]model.Config
}

//Konstruktor
func NewConfigInMemRepository() model.ConfigRepository {
	return ConfigInMemRepository{
		configs: make(map[string]model.Config),
	}
}

func (c ConfigInMemRepository) AddConfig(config model.Config) {
	key := fmt.Sprintf("%s/%d", config.Name, config.Version)
	c.configs[key] = config
	fmt.Println(c.configs)
}

func (c ConfigInMemRepository) GetConfig(name string, version int) (model.Config, error)  {
	key := fmt.Sprintf("%s/%d", name, version)
	config, ok := c.configs[key]
	// fmt.Println(config)
	if !ok {
		return model.Config{}, errors.New("config not found")
	}
	return config, nil
}