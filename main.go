package main

import (
	"net/http"
	"projekat/handlers"
	"projekat/model"
	"projekat/repository"
	"projekat/service"

	"github.com/gorilla/mux"
)

func main() {
	repo := repository.NewConfigInMemRepository()
	repoGroup := repository.NewConfigGroupInMemRepository()

	serviceConfig := service.NewConfigService(repo)
	serviceConfigGroup := service.NewConfigGroupService(repoGroup)

	params := make(map[string]string)
	params["username"] = "pera"
	params["port"] = "5432"
	config := model.Config{
		Name:       "viktorova",
		Version:    2,
		Parameters: params,
	}

	configs := []model.Config{}
	configs = append(configs, config)

	configGroup := model.ConfigGroup{
		Name:           "momirova",
		Version:        2,
		Configurations: configs,
	}

	serviceConfig.AddConfig(config)
	serviceConfigGroup.AddConfigGroup(configGroup)

	handlerConfig := handlers.NewConfigHandler(serviceConfig)
	handlerConfigGroup := handlers.NewConfigGroupHandler(serviceConfigGroup)

	router := mux.NewRouter()

	router.HandleFunc("/configs/{name}/{version}", handlerConfig.GetConfig).Methods("GET")
	router.HandleFunc("/configs", handlerConfig.AddConfig).Methods("POST")
	router.HandleFunc("/configs/{name}/{version}", handlerConfig.DeleteConfig).Methods("DELETE")

	router.HandleFunc("/configGroups/{name}/{version}", handlerConfigGroup.GetConfigGroup).Methods("GET")
	router.HandleFunc("/configGroups", handlerConfigGroup.AddConfigGroup).Methods("POST")
	router.HandleFunc("/configGroups/{name}/{version}", handlerConfigGroup.DeleteConfigGroup).Methods("DELETE")
	router.HandleFunc("/configGroups/{name}/{version}", handlerConfigGroup.AddConfigToGroup).Methods("POST")

	http.ListenAndServe("localhost:8000", router)
}
