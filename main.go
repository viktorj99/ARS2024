package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"projekat/handlers"
	"projekat/model"
	"projekat/repository"
	"projekat/service"
	"time"

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
	handlerConfigGroup := handlers.NewConfigGroupHandler(serviceConfigGroup, serviceConfig)

	router := mux.NewRouter()

	router.HandleFunc("/configs/{name}/{version}", handlerConfig.GetConfig).Methods("GET")
	router.HandleFunc("/configs", handlerConfig.AddConfig).Methods("POST")
	router.HandleFunc("/configs/{name}/{version}", handlerConfig.DeleteConfig).Methods("DELETE")

	router.HandleFunc("/configGroups/{name}/{version}", handlerConfigGroup.GetConfigGroup).Methods("GET")
	router.HandleFunc("/configGroups", handlerConfigGroup.AddConfigGroup).Methods("POST")
	router.HandleFunc("/configGroups/{name}/{version}", handlerConfigGroup.DeleteConfigGroup).Methods("DELETE")
	router.HandleFunc("/configGroups/{name}/{version}", handlerConfigGroup.AddConfigToGroup).Methods("POST")

	server := &http.Server{
		Addr:    "localhost:8000",
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-stop
	signal.Stop(stop)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
