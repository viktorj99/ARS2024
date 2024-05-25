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

	_, exists := c.service.GetConfig(config.Name, config.Version)
	if exists == nil {
		http.Error(writer, "Configuration with the given name and version already exists", http.StatusConflict)
		return
	}

	if len(config.Parameters) == 0 {
		fmt.Fprintf(writer, "Error: 'params' field is required and cannot be empty")
		return
	}

	if len(config.Labels) == 0 {
		fmt.Fprintf(writer, "Error: 'labels' field is required and cannot be empty")
		return
	}

	c.service.AddConfig(config)
	fmt.Fprintf(writer, "Received config: %+v", config)
}

// GET /configs/{name}/{version}
func (c ConfigHandler) GetConfig(writer http.ResponseWriter, request *http.Request) {
	// time.Sleep(10 * time.Second)
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

func (c ConfigHandler) DeleteConfig(writer http.ResponseWriter, request *http.Request) {
	name := mux.Vars(request)["name"]
	version := mux.Vars(request)["version"]

	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(writer, "Invalid version format", http.StatusBadRequest)
		return
	}

	err = c.service.DeleteConfig(name, versionInt)
	if err != nil {
		if err.Error() == "config not found" {
			http.Error(writer, err.Error(), http.StatusNotFound)
		} else {
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	response := map[string]string{"message": "Configuration successfully deleted"}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		http.Error(writer, "Failed to write response", http.StatusInternalServerError)
	}
}
