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
type ConfigGroupHandler struct {
	service service.ConfigGroupService
}

// Konstruktor
func NewConfigGroupHandler(service service.ConfigGroupService) ConfigGroupHandler {
	return ConfigGroupHandler{
		service: service,
	}
}

// GET /configs/{name}/{version}
func (c ConfigGroupHandler) GetConfigGroup(writer http.ResponseWriter, request *http.Request) {
	name := mux.Vars(request)["name"]
	version := mux.Vars(request)["version"]

	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	// pozovi servis metodu
	config, err := c.service.GetConfigGroup(name, versionInt)
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

func (c ConfigGroupHandler) AddConfigGroup(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	var configGroup model.ConfigGroup

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&configGroup)
	if err != nil {
		fmt.Fprintf(writer, "Error parsing JSON: %v", err)
		return
	}
	if configGroup.Name == "" {
		fmt.Fprintf(writer, "Error: 'name' field is required and cannot be empty")
		return
	}
	if configGroup.Version == 0 {
		fmt.Fprintf(writer, "Error: 'version' field is required and cannot be zero")
		return
	}
	if len(configGroup.Configurations) == 0 {
		fmt.Fprintf(writer, "Error: 'params' field is required and cannot be empty")
		return
	}

	c.service.AddConfigGroup(configGroup)
	fmt.Fprintf(writer, "Received config: %+v", configGroup)
}

func (c ConfigGroupHandler) DeleteConfigGroup(writer http.ResponseWriter, request *http.Request) {
	name := mux.Vars(request)["name"]
	version := mux.Vars(request)["version"]

	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	err = c.service.DeleteConfigGroup(name, versionInt)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}

	response := map[string]string{"message": "Configuration successfully deleted"}
	jsonResponse, _ := json.Marshal(response)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(jsonResponse)
}

func (c ConfigGroupHandler) AddConfigToGroup(writer http.ResponseWriter, request *http.Request) {

	defer request.Body.Close()

	name := mux.Vars(request)["name"]
	version := mux.Vars(request)["version"]

	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = c.service.GetConfigGroup(name, versionInt)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}

	var config model.Config

	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&config)
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

	c.service.AddConfigToGroup(name, versionInt, config)
	fmt.Fprintf(writer, "Received config: %+v", config)
}
