package main

import (
	"projekat/model"
	"projekat/repository"
	"projekat/service"
)

func main() {
	repo := repository.NewConfigInMemRepository()
	service := service.NewConfigService(repo)
	params := make(map[string]string)
	params["username"] = "pera"
	params["port"] = "5432"
	config := model.Config{
		Name:       "viktorova",
		Version:    2,
		Parameters: params,
	}
	config2 := model.Config{
		Name:       "viktorova",
		Version:    2,
		Parameters: params,
	}
	service.AddConfig(config)
	service.GetConfig(config.Name, config.Version)
	service.AddConfig(config2)
	service.GetConfig(config2.Name, config2.Version)

}
