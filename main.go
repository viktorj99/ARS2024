package main

import (
	"projekat/model"
	"projekat/repository"
)

func main() {
	repo := repository.NewConfigInMemRepository()

	params := make(map[string]string)
	params["username"] = "pera"
	params["port"] = "5432"
	config := model.Config{
		Name:    "viktorova",
		Version: 2,
		Parameters:  params,
	}
	
	repo.AddConfig(config)
	repo.GetConfig(config.Name, config.Version)
}
