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
	service := service.NewConfigService(repo)
	params := make(map[string]string)
	params["username"] = "pera"
	params["port"] = "5432"
	config := model.Config{
		Name:       "viktorova",
		Version:    2,
		Parameters: params,
	}

	service.AddConfig(config)
	handler := handlers.NewConfigHandler(service)

	router := mux.NewRouter()

	router.HandleFunc("/configs/{name}/{version}", handler.GetConfig).Methods("GET")
	router.HandleFunc("/configs", handler.AddConfig).Methods("POST")

	http.ListenAndServe("localhost:8000", router)
}
