package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"projekat/model"
	"projekat/service"
	"strconv"

	"github.com/gorilla/mux"
)

// Klasa
type ConfigHandler struct {
	service service.ConfigService
}

// Konstruktor
func NewConfigHandler(service service.ConfigService) ConfigHandler {
	return ConfigHandler{
		service: service,
	}
}

func (c ConfigHandler) AddConfig(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	var config model.Config

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Fprintf(writer, "Error parsing JSON: %v", err)
		return
	}
	if config.Name == "" {
		fmt.Fprintf(writer, "Error: 'name' field is required and cannot be empty")
		return
	}
	if config.Version == 0 {
		fmt.Fprintf(writer, "Error: 'version' field is required and cannot be zero")
		return
	}
	if len(config.Parameters) == 0 {
		fmt.Fprintf(writer, "Error: 'params' field is required and cannot be empty")
		return
	}

	c.service.AddConfig(config)
	fmt.Fprintf(writer, "Received config: %+v", config)
}

// GET /configs/{name}/{version}
func (c ConfigHandler) GetConfig(writer http.ResponseWriter, request *http.Request) {
	name := mux.Vars(request)["name"]
	version := mux.Vars(request)["version"]

	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	// pozovi servis metodu
	config, err := c.service.GetConfig(name, versionInt)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}

	// vrati odgovor
	response, err := json.Marshal(config)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content−Type", "application/json")
	writer.Write(response)
}
