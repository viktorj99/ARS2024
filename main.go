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

	_ "projekat/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/time/rate"
)

// @title Configuration API
// @version 1.0
// @description This is a sample server for a configuration service.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

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

	handlerConfig := handlers.NewConfigHandler(serviceConfig)
	handlerConfigGroup := handlers.NewConfigGroupHandler(serviceConfigGroup, serviceConfig)

	router := mux.NewRouter()

	limiter := rate.NewLimiter(0.167, 10)

	router.Handle("/configs/{name}/{version}", handlers.RateLimit(limiter, count(http.HandlerFunc(handlerConfig.GetConfig)))).Methods("GET")
	router.Handle("/configs", handlers.RateLimit(limiter, count(http.HandlerFunc(handlerConfig.AddConfig)))).Methods("POST")
	router.Handle("/configs/{name}/{version}", handlers.RateLimit(limiter, count(http.HandlerFunc(handlerConfig.DeleteConfig)))).Methods("DELETE")

	router.Handle("/configGroups/{name}/{version}", handlers.RateLimit(limiter, count(http.HandlerFunc(handlerConfigGroup.GetConfigGroup)))).Methods("GET")
	router.Handle("/configGroups", handlers.RateLimit(limiter, count(http.HandlerFunc(handlerConfigGroup.AddConfigGroup)))).Methods("POST")
	router.Handle("/configGroups/{name}/{version}", handlers.RateLimit(limiter, count(http.HandlerFunc(handlerConfigGroup.DeleteConfigGroup)))).Methods("DELETE")
	router.Handle("/configGroups/{name}/{version}", handlers.RateLimit(limiter, count(http.HandlerFunc(handlerConfigGroup.AddConfigToGroup)))).Methods("POST")
	router.Handle("/configGroups/{groupName}/{groupVersion}/{configName}/{configVersion}", handlers.RateLimit(limiter, count(http.HandlerFunc(handlerConfigGroup.DeleteConfigFromGroup)))).Methods("DELETE")
	router.Handle("/configGroups/{groupName}/{groupVersion}/{labels}", handlers.RateLimit(limiter, count(http.HandlerFunc(handlerConfigGroup.GetConfigsFromGroupByLabels)))).Methods("GET")
	router.Handle("/configGroups/{groupName}/{groupVersion}/{labels}", handlers.RateLimit(limiter, count(http.HandlerFunc(handlerConfigGroup.DeleteConfigsFromGroupByLabels)))).Methods("DELETE")

	// Swagger documentation route
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	router.Path("/metrics").Handler(metricsHandler())

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
