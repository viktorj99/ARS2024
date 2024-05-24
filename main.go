package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"projekat/handlers"
	"projekat/repository"
	"projekat/service"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

func main() {
	repo, err := repository.NewConfigConsulRepository()
	if err != nil {
		log.Fatalf("Failed to create Consul repository: %v", err)
	}

	repoGroup, err := repository.NewConfigGroupConsulRepository()
	if err != nil {
		log.Fatalf("Failed to create Consul repository: %v", err)
	}

	serviceConfig := service.NewConfigService(repo)
	serviceConfigGroup := service.NewConfigGroupService(repoGroup)

	// params := make(map[string]string)
	// params["username"] = "pera"
	// params["port"] = "5432"

	// labels := make(map[string]string)
	// labels["l1"] = "v1"
	// labels["l2"] = "v2"

	// config := model.Config{
	// 	Name:       "viktorova",
	// 	Version:    2,
	// 	Parameters: params,
	// 	Labels:     labels,
	// }

	// configs := []model.Config{}
	// configs = append(configs, config)

	// configGroup := model.ConfigGroup{
	// 	Name:           "momirova",
	// 	Version:        2,
	// 	Configurations: configs,
	// }

	// serviceConfig.AddConfig(config)
	// serviceConfigGroup.AddConfigGroup(configGroup)

	handlerConfig := handlers.NewConfigHandler(serviceConfig)
	handlerConfigGroup := handlers.NewConfigGroupHandler(serviceConfigGroup, serviceConfig)

	router := mux.NewRouter()

	limiter := rate.NewLimiter(0.167, 10)

	router.Handle("/configs/{name}/{version}", handlers.RateLimit(limiter, http.HandlerFunc(handlerConfig.GetConfig))).Methods("GET")
	router.Handle("/configs", handlers.RateLimit(limiter, http.HandlerFunc(handlerConfig.AddConfig))).Methods("POST")
	router.Handle("/configs/{name}/{version}", handlers.RateLimit(limiter, http.HandlerFunc(handlerConfig.DeleteConfig))).Methods("DELETE")

	router.Handle("/configGroups/{name}/{version}", handlers.RateLimit(limiter, http.HandlerFunc(handlerConfigGroup.GetConfigGroup))).Methods("GET")
	router.Handle("/configGroups", handlers.RateLimit(limiter, http.HandlerFunc(handlerConfigGroup.AddConfigGroup))).Methods("POST")
	router.Handle("/configGroups/{name}/{version}", handlers.RateLimit(limiter, http.HandlerFunc(handlerConfigGroup.DeleteConfigGroup))).Methods("DELETE")
	router.Handle("/configGroups/{name}/{version}", handlers.RateLimit(limiter, http.HandlerFunc(handlerConfigGroup.AddConfigToGroup))).Methods("POST")
	router.Handle("/configGroups/{groupName}/{groupVersion}/{configName}/{configVersion}", handlers.RateLimit(limiter, http.HandlerFunc(handlerConfigGroup.DeleteConfigFromGroup))).Methods("DELETE")
	router.Handle("/configGroups/{groupName}/{groupVersion}/{labels}", handlers.RateLimit(limiter, http.HandlerFunc(handlerConfigGroup.GetConfigsFromGroupByLabels))).Methods("GET")
	router.Handle("/configGroups/{groupName}/{groupVersion}/{labels}", handlers.RateLimit(limiter, http.HandlerFunc(handlerConfigGroup.DeleteConfigsFromGroupByLabels))).Methods("DELETE")

	server := &http.Server{
		Addr:    "0.0.0.0:8080",
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
